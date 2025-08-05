package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) GetImages(ctx context.Context, args v1.GetImageArgs) ([]v1.Image, error) {
	return nil, v1.ErrNotImplemented
}
