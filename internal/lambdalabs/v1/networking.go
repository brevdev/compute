package v1

import (
	"context"
	"fmt"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// AddFirewallRulesToInstance adds firewall rules to an instance
func (c *LambdaLabsClient) AddFirewallRulesToInstance(ctx context.Context, args v1.AddFirewallRulesToInstanceArgs) error {
	// TODO: Implement Lambda Labs firewall rule addition
	// This would typically involve:
	// 1. Validating the firewall rules
	// 2. Calling Lambda Labs API to add the rules
	// 3. Waiting for the operation to complete

	return fmt.Errorf("not implemented")
}

// RevokeSecurityGroupRules revokes security group rules from an instance
func (c *LambdaLabsClient) RevokeSecurityGroupRules(ctx context.Context, args v1.RevokeSecurityGroupRuleArgs) error {
	// TODO: Implement Lambda Labs security group rule revocation
	return fmt.Errorf("not implemented")
}
