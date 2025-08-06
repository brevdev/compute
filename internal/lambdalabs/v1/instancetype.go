package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/units"
	"github.com/bojanz/currency"
	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/compute/pkg/v1"
)

// GetInstanceTypes retrieves available instance types from Lambda Labs
// Supported via: GET /api/v1/instance-types
func (c *LambdaLabsClient) GetInstanceTypes(ctx context.Context, args v1.GetInstanceTypeArgs) ([]v1.InstanceType, error) {
	resp, httpResp, err := c.client.DefaultAPI.InstanceTypes(c.makeAuthContext(ctx)).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, handleLLErrToCloudErr(err)
	}

	var instanceTypes []v1.InstanceType

	for instanceTypeName, instanceTypeData := range resp.Data {
		if strings.Contains(instanceTypeName, "gh") {
			continue
		}

		if len(args.InstanceTypes) > 0 && !contains(args.InstanceTypes, instanceTypeName) {
			continue
		}

		availableRegions := make(map[string]bool)
		for _, region := range instanceTypeData.RegionsWithCapacityAvailable {
			availableRegions[region.Name] = true
		}

		for _, region := range instanceTypeData.RegionsWithCapacityAvailable {
		if len(args.Locations) > 0 && !args.Locations.IsAll() && !containsLocation(args.Locations, region.Name) {
			continue
		}

		v1InstanceType, err := convertLambdaLabsInstanceTypeToV1InstanceType(region.Name, instanceTypeData.InstanceType, true)
		if err != nil {
			return nil, fmt.Errorf("failed to convert instance type: %w", err)
		}
		instanceTypes = append(instanceTypes, v1InstanceType)
	}
	}

	return instanceTypes, nil
}

// GetInstanceTypePollTime returns the polling interval for instance types
func (c *LambdaLabsClient) GetInstanceTypePollTime() time.Duration {
	// TODO: Configure appropriate polling time for Lambda Labs
	return 5 * time.Minute
}

// GetLocations retrieves available locations from Lambda Labs
// UNSUPPORTED: No location listing endpoints found in Lambda Labs API
func convertLambdaLabsInstanceTypeToV1InstanceType(location string, llInstanceType openapi.InstanceType, isAvailable bool) (v1.InstanceType, error) {
	var gpus []v1.GPU
	if !strings.Contains(llInstanceType.Description, "CPU") {
		gpu := parseGPUFromDescription(llInstanceType.Description)
		gpus = append(gpus, gpu)
	}

	amount, err := currency.NewAmountFromInt64(int64(llInstanceType.PriceCentsPerHour), "USD")
	if err != nil {
		return v1.InstanceType{}, fmt.Errorf("failed to create price amount: %w", err)
	}

	instanceType := v1.InstanceType{
		Location:      location,
		Type:          llInstanceType.Name,
		SupportedGPUs: gpus,
		SupportedStorage: []v1.Storage{
			{
				Type: "ssd",
				Size: units.GiB * units.Base2Bytes(llInstanceType.Specs.StorageGib),
			},
		},
		Memory:                 units.GiB * units.Base2Bytes(llInstanceType.Specs.MemoryGib),
		VCPU:                   llInstanceType.Specs.Vcpus,
		SupportedArchitectures: []string{"x86_64"},
		Stoppable:              false,
		Rebootable:             true,
		IsAvailable:            isAvailable,
		BasePrice:              &amount,
		Provider:               "lambdalabs",
	}

	instanceType.ID = v1.InstanceTypeID(fmt.Sprintf("lambdalabs-%s-%s", location, llInstanceType.Name))

	return instanceType, nil
}

const lambdaLocationsData = `[
    {"location_name": "us-west-1", "description": "California, USA", "country": "USA"},
    {"location_name": "us-west-2", "description": "Arizona, USA", "country": "USA"},
    {"location_name": "us-west-3", "description": "Utah, USA", "country": "USA"},
    {"location_name": "us-south-1", "description": "Texas, USA", "country": "USA"},
    {"location_name": "us-east-1", "description": "Virginia, USA", "country": "USA"},
    {"location_name": "us-midwest-1", "description": "Illinois, USA", "country": "USA"},
    {"location_name": "australia-southeast-1", "description": "Australia", "country": "AUS"},
    {"location_name": "europe-central-1", "description": "Germany", "country": "DEU"},
    {"location_name": "asia-south-1", "description": "India", "country": "IND"},
    {"location_name": "me-west-1", "description": "Israel", "country": "ISR"},
    {"location_name": "europe-south-1", "description": "Italy", "country": "ITA"},
    {"location_name": "asia-northeast-1", "description": "Osaka, Japan", "country": "JPN"},
    {"location_name": "asia-northeast-2", "description": "Tokyo, Japan", "country": "JPN"},
    {"location_name": "us-east-3", "description": "Washington D.C, USA", "country": "USA"},
    {"location_name": "us-east-2", "description": "Washington D.C, USA", "country": "USA"},
    {"location_name": "australia-east-1", "description": "Sydney, Australia", "country": "AUS"},
    {"location_name": "us-south-3", "description": "Central Texas, USA", "country": "USA"},
    {"location_name": "us-south-2", "description": "North Texas, USA", "country": "USA"}
]`

type LambdaLocation struct {
	LocationName string `json:"location_name"`
	Description  string `json:"description"`
	Country      string `json:"country"`
}

func (c *LambdaLabsClient) GetLocations(_ context.Context, _ v1.GetLocationsArgs) ([]v1.Location, error) {
	var regionData []LambdaLocation
	if err := json.Unmarshal([]byte(lambdaLocationsData), &regionData); err != nil {
		return nil, fmt.Errorf("failed to parse location data: %w", err)
	}

	locations := make([]v1.Location, 0, len(regionData))
	for _, region := range regionData {
		locations = append(locations, v1.Location{
			Name:        region.LocationName,
			Description: region.Description,
			Available:   true,
			Country:     region.Country,
		})
	}

	return locations, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsLocation(locations v1.LocationsFilter, location string) bool {
	for _, loc := range locations {
		if loc == location {
			return true
		}
	}
	return false
}
