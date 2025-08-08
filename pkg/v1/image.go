package v1

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/brevdev/cloud/pkg/ssh"
)

type CloudMachineImage interface {
	GetImages(ctx context.Context, args GetImageArgs) ([]Image, error)
}

type GetImageArgs struct {
	Owners        []string // self, amazon, aws-marketplace, <account id>, project id for GCP
	Architectures []string // i386, x86_64, arm64
	NameFilters   []string // name of the image (wildcard permitted)
	ImageIDs      []string
}

type Image struct {
	ID           string
	Architecture string
	Description  string
	Name         string
	CreatedAt    time.Time
}

func ValidateInstanceImage(ctx context.Context, instance Instance, privateKey string) error {
	// First ensure the instance is running and SSH accessible
	sshUser := instance.SSHUser
	sshPort := instance.SSHPort
	publicIP := instance.PublicIP

	// Validate that we have the required SSH connection details
	if sshUser == "" {
		return fmt.Errorf("SSH user is not set for instance %s", instance.CloudID)
	}
	if sshPort == 0 {
		return fmt.Errorf("SSH port is not set for instance %s", instance.CloudID)
	}
	if publicIP == "" {
		return fmt.Errorf("public IP is not available for instance %s", instance.CloudID)
	}

	// Connect to the instance via SSH
	sshClient, err := ssh.ConnectToHost(ctx, ssh.ConnectionConfig{
		User:     sshUser,
		HostPort: fmt.Sprintf("%s:%d", publicIP, sshPort),
		PrivKey:  privateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to instance via SSH: %w", err)
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			// Log close error but don't return it as it's not the primary error
			fmt.Printf("warning: failed to close SSH connection: %v\n", closeErr)
		}
	}()

	// Check 1: Verify x86_64 architecture
	stdout, stderr, err := sshClient.RunCommand(ctx, "uname -m")
	if err != nil {
		return fmt.Errorf("failed to check architecture: %w, stdout: %s, stderr: %s", err, stdout, stderr)
	}
	if !strings.Contains(strings.TrimSpace(stdout), "x86_64") {
		return fmt.Errorf("expected x86_64 architecture, got: %s", strings.TrimSpace(stdout))
	}

	// Check 2: Verify Ubuntu 20.04 or 22.04
	stdout, stderr, err = sshClient.RunCommand(ctx, "cat /etc/os-release | grep PRETTY_NAME")
	if err != nil {
		return fmt.Errorf("failed to check OS version: %w, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	parts := strings.Split(strings.TrimSpace(stdout), "=")
	if len(parts) != 2 {
		return fmt.Errorf("error: os pretty name not in format PRETTY_NAME=\"Ubuntu\": %s", stdout)
	}

	// Remove quotes from the value
	osVersion := strings.Trim(parts[1], "\"")
	ubuntuRegex := regexp.MustCompile(`Ubuntu 20\.04|22\.04`)
	if !ubuntuRegex.MatchString(osVersion) {
		return fmt.Errorf("expected Ubuntu 20.04 or 22.04, got: %s", osVersion)
	}

	// Check 3: Verify home directory
	stdout, stderr, err = sshClient.RunCommand(ctx, "cd ~ && pwd")
	if err != nil {
		return fmt.Errorf("failed to check home directory: %w, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	homeDir := strings.TrimSpace(stdout)
	if sshUser == "ubuntu" {
		if !strings.Contains(homeDir, "/home/ubuntu") {
			return fmt.Errorf("expected ubuntu user home directory to contain /home/ubuntu, got: %s", homeDir)
		}
	} else {
		if !strings.Contains(homeDir, "/root") {
			return fmt.Errorf("expected non-ubuntu user home directory to contain /root, got: %s", homeDir)
		}
	}

	fmt.Printf("Instance image validation passed for %s: architecture=%s, os=%s, home=%s\n",
		instance.CloudID, "x86_64", osVersion, homeDir)

	return nil
}
