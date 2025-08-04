package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// GetCapabilities returns the capabilities of Lambda Labs
func (c *LambdaLabsClient) GetCapabilities(ctx context.Context) (v1.Capabilities, error) {
	// TODO: Implement Lambda Labs capabilities
	// This would typically involve:
	// 1. Determining what features Lambda Labs supports
	// 2. Returning a Capabilities struct with the supported features

	capabilities := v1.Capabilities{
		// TODO: Fill in actual Lambda Labs capabilities
		// Example capabilities that Lambda Labs might support:
		// - GPU instances
		// - Spot instances
		// - Custom images
		// - Firewall rules
		// - Volume resizing
		// - Instance type changing
		// - Tagging
		// - etc.
	}

	return capabilities, nil
}
