package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) GetImages(ctx context.Context, args v1.GetImageArgs) ([]v1.Image, error) {
	return nil, v1.ErrNotImplemented
}
