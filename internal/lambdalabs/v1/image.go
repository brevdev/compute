package v1

import (
	"context"
	"fmt"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// GetImages retrieves available images from Lambda Labs
func (c *LambdaLabsClient) GetImages(ctx context.Context, args v1.GetImageArgs) ([]v1.Image, error) {
	// TODO: Implement Lambda Labs image retrieval
	// This would typically involve:
	// 1. Calling Lambda Labs API to get available images
	// 2. Filtering based on the provided arguments
	// 3. Converting to the standard Image format

	return nil, fmt.Errorf("not implemented")
}
