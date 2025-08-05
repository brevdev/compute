package v1

import (
	"context"
	"time"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) GetInstanceTypes(ctx context.Context, args v1.GetInstanceTypeArgs) ([]v1.InstanceType, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) GetInstanceTypePollTime() time.Duration {
	return 5 * time.Minute
}

func (c *FluidStackClient) GetLocations(ctx context.Context, args v1.GetLocationsArgs) ([]v1.Location, error) {
	return nil, v1.ErrNotImplemented
}
