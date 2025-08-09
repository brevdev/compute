package v1

import (
	"context"
	"fmt"

	v1 "github.com/brevdev/cloud/pkg/v1"
)

func (c *NebiusClient) AddFirewallRulesToInstance(ctx context.Context, args v1.AddFirewallRulesToInstanceArgs) error {
	securityGroupID, err := c.getOrCreateSecurityGroupForInstance(ctx, args.InstanceID)
	if err != nil {
		return fmt.Errorf("failed to get or create security group for instance %s: %w", args.InstanceID, err)
	}

	err = c.addFirewallRulesToSecurityGroup(ctx, securityGroupID, args.FirewallRules)
	if err != nil {
		return fmt.Errorf("failed to add firewall rules to security group %s: %w", securityGroupID, err)
	}

	return nil
}

func (c *NebiusClient) RevokeSecurityGroupRules(ctx context.Context, args v1.RevokeSecurityGroupRuleArgs) error {
	securityGroupID, err := c.getSecurityGroupForInstance(ctx, args.InstanceID)
	if err != nil {
		return fmt.Errorf("failed to get security group for instance %s: %w", args.InstanceID, err)
	}

	err = c.removeSecurityGroupRules(ctx, securityGroupID, args.SecurityGroupRuleIDs)
	if err != nil {
		return fmt.Errorf("failed to remove security group rules from %s: %w", securityGroupID, err)
	}

	return nil
}

func (c *NebiusClient) getOrCreateSecurityGroupForInstance(_ context.Context, instanceID v1.CloudProviderInstanceID) (string, error) {
	clusterID := c.getClusterIDFromInstance(instanceID)
	_ = fmt.Sprintf("brev-cluster-%s", clusterID)

	return "", fmt.Errorf("security group management not yet implemented - need to use Nebius VPC service")
}

func (c *NebiusClient) getSecurityGroupForInstance(_ context.Context, _ v1.CloudProviderInstanceID) (string, error) {
	return "", fmt.Errorf("security group lookup not yet implemented - need to use Nebius VPC service")
}

func (c *NebiusClient) addFirewallRulesToSecurityGroup(_ context.Context, _ string, _ v1.FirewallRules) error {
	return fmt.Errorf("firewall rule addition not yet implemented - need to use Nebius VPC service")
}

func (c *NebiusClient) removeSecurityGroupRules(_ context.Context, _ string, _ []string) error {
	return fmt.Errorf("security group rule removal not yet implemented - need to use Nebius VPC service")
}

func (c *NebiusClient) getClusterIDFromInstance(_ v1.CloudProviderInstanceID) string {
	return "default"
}
