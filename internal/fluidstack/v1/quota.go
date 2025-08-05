package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) GetInstanceTypeQuotas(_ context.Context, _ v1.GetInstanceTypeQuotasArgs) (v1.Quota, error) {
	return v1.Quota{}, v1.ErrNotImplemented
}
