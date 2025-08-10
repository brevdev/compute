package v1

import (
	"context"
	"fmt"

	"github.com/brevdev/cloud/pkg/ssh"
	"github.com/google/uuid"
)

type CloudModifyFirewall interface {
	AddFirewallRulesToInstance(ctx context.Context, args AddFirewallRulesToInstanceArgs) error
	RevokeSecurityGroupRules(ctx context.Context, args RevokeSecurityGroupRuleArgs) error
}

type AddFirewallRulesToInstanceArgs struct {
	InstanceID    CloudProviderInstanceID
	FirewallRules FirewallRules
}

type RevokeSecurityGroupRuleArgs struct {
	InstanceID           CloudProviderInstanceID
	SecurityGroupRuleIDs []string
}

type FirewallRules struct {
	IngressRules []FirewallRule
	EgressRules  []FirewallRule
}

type FirewallRule struct {
	ID       string // ignored when creating a new rule
	FromPort int32
	ToPort   int32
	IPRanges []string
}

type PortMapping struct {
	FromPort int
	ToPort   int
}

func ValidateInboundPortRestriction(ctx context.Context, client CloudInstanceReader, instance *Instance, privateKey string) error {
	var err error
	instance, err = WaitForInstanceLifecycleStatus(ctx, client, instance, LifecycleStatusRunning, PendingToRunningTimeout)
	if err != nil {
		return err
	}

	if instance.SSHUser == "" {
		return fmt.Errorf("SSH user is not set for instance %s", instance.CloudID)
	}
	if instance.SSHPort == 0 {
		return fmt.Errorf("SSH port is not set for instance %s", instance.CloudID)
	}
	if instance.PublicIP == "" {
		return fmt.Errorf("public IP is not available for instance %s", instance.CloudID)
	}

	sshClient, err := ssh.ConnectToHost(ctx, ssh.ConnectionConfig{
		User:     instance.SSHUser,
		HostPort: fmt.Sprintf("%s:%d", instance.PublicIP, instance.SSHPort),
		PrivKey:  privateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to instance via SSH: %w", err)
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close SSH connection: %v\n", closeErr)
		}
	}()

	portsToCheck := []int{21, 23, 25, 53, 80, 443, 993, 995, 3389, 5432, 3306}

	for _, port := range portsToCheck {
		cmd := fmt.Sprintf("timeout 5 nc -z %s %d", instance.PublicIP, port)
		stdout, stderr, err := sshClient.RunCommand(ctx, cmd)

		if err == nil {
			return fmt.Errorf("security violation: port %d is accessible from external sources, stdout: %s, stderr: %s", port, stdout, stderr)
		}

		fmt.Printf("Port %d properly blocked (expected): %s\n", port, stderr)
	}

	cmd := fmt.Sprintf("timeout 5 nc -z %s %d", instance.PublicIP, instance.SSHPort)
	stdout, stderr, err := sshClient.RunCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("SSH port %d should be accessible but is not: %w, stdout: %s, stderr: %s", instance.SSHPort, err, stdout, stderr)
	}

	fmt.Printf("Inbound port validation passed: only SSH port %d is accessible\n", instance.SSHPort)
	return nil
}

func ValidateEastWestConnectivity(ctx context.Context, client CloudCreateTerminateInstance, attrs CreateInstanceAttrs, privateKey string) error {
	instance1, instance2, err := createTestInstances(ctx, client, attrs)
	if err != nil {
		return err
	}

	defer cleanupInstances(ctx, client, instance1, instance2)

	err = waitForInstancesReady(ctx, client, instance1, instance2, privateKey)
	if err != nil {
		return err
	}

	return testConnectivity(ctx, instance1, instance2, privateKey)
}

func createTestInstances(ctx context.Context, client CloudCreateTerminateInstance, attrs CreateInstanceAttrs) (*Instance, *Instance, error) {
	attrs1 := attrs
	attrs1.RefID = uuid.New().String()
	attrs1.Name = fmt.Sprintf("%s-east", attrs.Name)

	instance1, err := client.CreateInstance(ctx, attrs1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create first instance: %w", err)
	}

	attrs2 := attrs
	attrs2.RefID = uuid.New().String()
	attrs2.Name = fmt.Sprintf("%s-west", attrs.Name)

	instance2, err := client.CreateInstance(ctx, attrs2)
	if err != nil {
		return instance1, nil, fmt.Errorf("failed to create second instance: %w", err)
	}

	return instance1, instance2, nil
}

func cleanupInstances(ctx context.Context, client CloudCreateTerminateInstance, instance1, instance2 *Instance) {
	if instance1 != nil {
		if termErr := client.TerminateInstance(ctx, instance1.CloudID); termErr != nil {
			fmt.Printf("warning: failed to terminate first instance %s: %v\n", instance1.CloudID, termErr)
		}
	}
	if instance2 != nil {
		if termErr := client.TerminateInstance(ctx, instance2.CloudID); termErr != nil {
			fmt.Printf("warning: failed to terminate second instance %s: %v\n", instance2.CloudID, termErr)
		}
	}
}

func waitForInstancesReady(ctx context.Context, client CloudCreateTerminateInstance, instance1, instance2 *Instance, privateKey string) error {
	var err error
	instance1, err = WaitForInstanceLifecycleStatus(ctx, client, instance1, LifecycleStatusRunning, PendingToRunningTimeout)
	if err != nil {
		return fmt.Errorf("first instance failed to reach running state: %w", err)
	}

	instance2, err = WaitForInstanceLifecycleStatus(ctx, client, instance2, LifecycleStatusRunning, PendingToRunningTimeout)
	if err != nil {
		return fmt.Errorf("second instance failed to reach running state: %w", err)
	}

	err = ssh.WaitForSSH(ctx, ssh.ConnectionConfig{
		User:     instance1.SSHUser,
		HostPort: fmt.Sprintf("%s:%d", instance1.PublicIP, instance1.SSHPort),
		PrivKey:  privateKey,
	}, ssh.WaitForSSHOptions{
		Timeout: RunningSSHTimeout,
	})
	if err != nil {
		return fmt.Errorf("SSH not available on first instance: %w", err)
	}

	err = ssh.WaitForSSH(ctx, ssh.ConnectionConfig{
		User:     instance2.SSHUser,
		HostPort: fmt.Sprintf("%s:%d", instance2.PublicIP, instance2.SSHPort),
		PrivKey:  privateKey,
	}, ssh.WaitForSSHOptions{
		Timeout: RunningSSHTimeout,
	})
	if err != nil {
		return fmt.Errorf("SSH not available on second instance: %w", err)
	}

	return nil
}

func testConnectivity(ctx context.Context, instance1, instance2 *Instance, privateKey string) error {
	sshClient1, err := ssh.ConnectToHost(ctx, ssh.ConnectionConfig{
		User:     instance1.SSHUser,
		HostPort: fmt.Sprintf("%s:%d", instance1.PublicIP, instance1.SSHPort),
		PrivKey:  privateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to first instance via SSH: %w", err)
	}
	defer func() {
		if closeErr := sshClient1.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close SSH connection to first instance: %v\n", closeErr)
		}
	}()

	pingCmd := fmt.Sprintf("ping -c 3 -W 5 %s", instance2.PrivateIP)
	stdout, stderr, err := sshClient1.RunCommand(ctx, pingCmd)
	if err != nil {
		return fmt.Errorf("ping from instance1 to instance2 failed: %w, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	sshTestCmd := fmt.Sprintf("timeout 10 nc -z %s %d", instance2.PrivateIP, instance2.SSHPort)
	stdout, stderr, err = sshClient1.RunCommand(ctx, sshTestCmd)
	if err != nil {
		fmt.Printf("SSH port connectivity test between instances failed (may be expected): %s, stderr: %s\n", stdout, stderr)
	} else {
		fmt.Printf("SSH port connectivity between instances successful: %s\n", stdout)
	}

	fmt.Printf("East-west connectivity validation passed: instance1 (%s) can communicate with instance2 (%s)\n",
		instance1.CloudID, instance2.CloudID)
	return nil
}
