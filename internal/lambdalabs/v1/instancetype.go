package v1

import (
	"context"
	"time"

	"github.com/alecthomas/units"
	v1 "github.com/brevdev/compute/pkg/v1"
)

// GetInstanceTypes retrieves available instance types from Lambda Labs
func (c *LambdaLabsClient) GetInstanceTypes(ctx context.Context, args v1.GetInstanceTypeArgs) ([]v1.InstanceType, error) {
	// TODO: Implement Lambda Labs instance type retrieval
	// This would typically involve:
	// 1. Calling Lambda Labs API to get available instance types
	// 2. Filtering based on the provided arguments
	// 3. Converting to the standard InstanceType format

	// Example stub implementation
	instanceTypes := []v1.InstanceType{
		{
			ID:                     v1.InstanceTypeID("gpu_1x_a10"),
			Location:               "us-east-1",
			AvailableAzs:           []string{"us-east-1a", "us-east-1b"},
			SubLocation:            "us-east-1a",
			Type:                   "gpu_1x_a10",
			SupportedGPUs:          []v1.GPU{{Count: 1, Memory: 24 * units.GiB, Manufacturer: "NVIDIA", Name: "A10", Type: "A10"}},
			SupportedStorage:       []v1.Storage{{Type: "ssd", Size: 100 * units.GiB}},
			ElasticRootVolume:      true,
			Memory:                 24 * units.GiB,
			VCPU:                   4,
			SupportedArchitectures: []string{"x86_64"},
			IsAvailable:            true,
			BasePrice:              nil, // TODO: Get actual pricing using currency.New
			Provider:               "lambdalabs",
		},
	}

	return instanceTypes, nil
}

// GetInstanceTypePollTime returns the polling interval for instance types
func (c *LambdaLabsClient) GetInstanceTypePollTime() time.Duration {
	// TODO: Configure appropriate polling time for Lambda Labs
	return 5 * time.Minute
}

// GetLocations retrieves available locations from Lambda Labs
func (c *LambdaLabsClient) GetLocations(ctx context.Context, args v1.GetLocationsArgs) ([]v1.Location, error) {
	// TODO: Implement Lambda Labs location retrieval
	locations := []v1.Location{
		{
			Name:        "us-east-1",
			Description: "US East (N. Virginia)",
			Available:   true,
			Endpoint:    "https://cloud.lambdalabs.com/api/v1",
			Priority:    1,
			Country:     "USA",
		},
		{
			Name:        "us-west-1",
			Description: "US West (N. California)",
			Available:   true,
			Endpoint:    "https://cloud.lambdalabs.com/api/v1",
			Priority:    2,
			Country:     "USA",
		},
	}

	return locations, nil
}
