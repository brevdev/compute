package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) AddFirewallRulesToInstance(ctx context.Context, args v1.AddFirewallRulesToInstanceArgs) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) RevokeSecurityGroupRules(ctx context.Context, args v1.RevokeSecurityGroupRuleArgs) error {
	return v1.ErrNotImplemented
}
