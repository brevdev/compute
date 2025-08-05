package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) GetImages(_ context.Context, _ v1.GetImageArgs) ([]v1.Image, error) {
	return nil, v1.ErrNotImplemented
}
