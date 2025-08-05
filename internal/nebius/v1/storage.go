package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) ResizeInstanceVolume(ctx context.Context, args v1.ResizeInstanceVolumeArgs) error {
	return v1.ErrNotImplemented
}
