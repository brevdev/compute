package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) UpdateInstanceTags(_ context.Context, args v1.UpdateInstanceTagsArgs) error {
	return v1.ErrNotImplemented
}
