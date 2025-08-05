package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) GetLocations(ctx context.Context, args v1.GetLocationsArgs) ([]v1.Location, error) {
	return nil, v1.ErrNotImplemented
}
