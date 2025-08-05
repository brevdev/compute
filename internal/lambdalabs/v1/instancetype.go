package v1

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
		return nil, fmt.Errorf("failed to get instance types: %w", err)
	}

	var instanceTypes []v1.InstanceType
	for _, llInstanceTypeData := range resp.Data {
		for _, region := range llInstanceTypeData.RegionsWithCapacityAvailable {
			instanceType, err := convertLambdaLabsInstanceTypeToV1InstanceType(
				region.Name,
				llInstanceTypeData.InstanceType,
				true,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to convert instance type: %w", err)
			}
			instanceTypes = append(instanceTypes, instanceType)
		}
	}

	if len(args.Locations) > 0 && !args.Locations.IsAll() {
		filtered := make([]v1.InstanceType, 0)
		for _, it := range instanceTypes {
			for _, loc := range args.Locations {
				if it.Location == loc {
					filtered = append(filtered, it)
					break
				}
			}
		}
		instanceTypes = filtered
	}

	if len(args.InstanceTypes) > 0 {
		filtered := make([]v1.InstanceType, 0)
		for _, it := range instanceTypes {
			for _, itName := range args.InstanceTypes {
				if it.Type == itName {
					filtered = append(filtered, it)
					break
				}
			}
		}
		instanceTypes = filtered
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

func parseGPUFromDescription(description string) v1.GPU {
	countRegex := regexp.MustCompile(`(\d+)x`)
	memoryRegex := regexp.MustCompile(`(\d+) GB`)
	nameRegex := regexp.MustCompile(`x (.*?) \(`)

	var gpu v1.GPU

	if matches := countRegex.FindStringSubmatch(description); len(matches) > 1 {
		if count, err := strconv.ParseInt(matches[1], 10, 32); err == nil {
			gpu.Count = int32(count)
		}
	}

	if matches := memoryRegex.FindStringSubmatch(description); len(matches) > 1 {
		if memory, err := strconv.Atoi(matches[1]); err == nil {
			gpu.Memory = units.GiB * units.Base2Bytes(memory)
		}
	}

	if matches := nameRegex.FindStringSubmatch(description); len(matches) > 1 {
		gpu.Name = strings.TrimSpace(matches[1])
		gpu.Type = gpu.Name
	}

	gpu.Manufacturer = "NVIDIA"

	return gpu
}

func (c *LambdaLabsClient) GetLocations(ctx context.Context, _ v1.GetLocationsArgs) ([]v1.Location, error) {
	resp, httpResp, err := c.client.DefaultAPI.InstanceTypes(c.makeAuthContext(ctx)).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get instance types: %w", err)
	}

	locationMap := make(map[string]bool)
	for _, llInstanceTypeData := range resp.Data {
		for _, region := range llInstanceTypeData.RegionsWithCapacityAvailable {
			locationMap[region.Name] = true
		}
	}

	var locations []v1.Location
	for locationName := range locationMap {
		locations = append(locations, v1.Location{
			Name:        locationName,
			Description: fmt.Sprintf("Lambda Labs region: %s", locationName),
			Available:   true,
			Country:     "USA",
		})
	}

	return locations, nil
}
