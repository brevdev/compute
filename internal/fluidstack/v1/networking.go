package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) AddFirewallRulesToInstance(ctx context.Context, args v1.AddFirewallRulesToInstanceArgs) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) RevokeSecurityGroupRules(ctx context.Context, args v1.RevokeSecurityGroupRuleArgs) error {
	return v1.ErrNotImplemented
}
