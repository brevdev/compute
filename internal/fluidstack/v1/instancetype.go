package v1

import (
	"context"
	"time"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) GetInstanceTypes(_ context.Context, _ v1.GetInstanceTypeArgs) ([]v1.InstanceType, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) GetInstanceTypePollTime() time.Duration {
	return 5 * time.Minute
}

func (c *FluidStackClient) GetLocations(_ context.Context, _ v1.GetLocationsArgs) ([]v1.Location, error) {
	return nil, v1.ErrNotImplemented
}
