package v1

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/alecthomas/units"
	"github.com/bojanz/currency"
)

type InstanceTypeID string

type InstanceType struct {
	ID                              InstanceTypeID // this id should be unique across all regions and stable
	Location                        string
	AvailableAzs                    []string
	SubLocation                     string
	Type                            string
	SupportedGPUs                   []GPU
	SupportedStorage                []Storage
	ElasticRootVolume               bool
	SupportedUsageClasses           []string
	Memory                          units.Base2Bytes
	MaximumNetworkInterfaces        int32
	NetworkPerformance              string
	SupportedNumCores               []int32
	DefaultCores                    int32
	VCPU                            int32
	SupportedArchitectures          []string
	ClockSpeedInGhz                 float64
	Quota                           InstanceTypeQuota
	Stoppable                       bool
	Rebootable                      bool
	VariablePrice                   bool
	Preemptible                     bool
	IsAvailable                     bool
	BasePrice                       *currency.Amount
	SubLocationTypeChangeable       bool
	IsContainer                     bool
	UserPrivilegeEscalationDisabled bool
	NotPrivileged                   bool
	EstimatedDeployTime             *time.Duration
	Provider                        string
	CanModifyFirewallRules          bool
}

type GPU struct {
	Count          int32
	Memory         units.Base2Bytes
	MemoryDetails  string
	NetworkDetails string
	Manufacturer   string
	Name           string
	Type           string
}

type InstanceTypeQuota struct {
	OnDemand Quota
	Spot     Quota
	Reserved Quota
}

type CloudInstanceType interface {
	GetInstanceTypes(ctx context.Context, args GetInstanceTypeArgs) ([]InstanceType, error)
	GetInstanceTypePollTime() time.Duration
	CloudLocation
}

type GetInstanceTypeArgs struct {
	Locations              LocationsFilter
	SupportedArchitectures []string
	InstanceTypes          []string
}

// ValidateGetInstanceTypes validates that the GetInstanceTypes functionality works correctly
// by testing that filtering by specific instance types returns the expected results
func ValidateGetInstanceTypes(ctx context.Context, client CloudInstanceType) error {
	// Get all instance types first
	allTypes, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{})
	if err != nil {
		return fmt.Errorf("failed to get all instance types: %w", err)
	}

	if len(allTypes) == 0 {
		return errors.New("no instance types available for validation")
	}

	// Test 1: Deterministic results - multiple calls should return the same results
	allTypes2, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{})
	if err != nil {
		return fmt.Errorf("failed to get all instance types on second call: %w", err)
	}

	// Remove volatile fields for comparison
	normalizedTypes1 := normalizeInstanceTypes(allTypes)
	normalizedTypes2 := normalizeInstanceTypes(allTypes2)

	if !reflect.DeepEqual(normalizedTypes1, normalizedTypes2) {
		return fmt.Errorf("instance types are not deterministic between calls")
	}

	// Test 2: ID stability and uniqueness
	idMap := make(map[InstanceTypeID]InstanceType)
	for _, instanceType := range allTypes {
		if existing, exists := idMap[instanceType.ID]; exists {
			return fmt.Errorf("duplicate instance type ID found: %s (types: %s, %s)",
				instanceType.ID, existing.Type, instanceType.Type)
		}
		idMap[instanceType.ID] = instanceType
	}

	// Test 3: Filtering by instance type name
	firstType := allTypes[0]
	filteredTypes, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{
		InstanceTypes: []string{firstType.Type},
	})
	if err != nil {
		return fmt.Errorf("failed to get filtered instance types: %w", err)
	}

	if len(filteredTypes) == 0 {
		return fmt.Errorf("no instance types returned when filtering by type: %s", firstType.Type)
	}

	// Compare the first type with the filtered result, ignoring fields that may vary
	// between different calls or implementations
	expectedType := firstType
	expectedType.ID = ""
	expectedType.SubLocation = ""
	expectedType.AvailableAzs = nil

	actualType := filteredTypes[0]
	actualType.ID = ""
	actualType.SubLocation = ""
	actualType.AvailableAzs = nil

	// Use reflection to compare the structs
	if !reflect.DeepEqual(expectedType, actualType) {
		return fmt.Errorf("filtered instance type does not match expected type: expected %+v, got %+v", expectedType, actualType)
	}

	return nil
}

// ValidateRegionalInstanceTypes validates that regional filtering works correctly
// by comparing regional results with all-region results using CloudLocation capabilities
func ValidateRegionalInstanceTypes(ctx context.Context, client CloudInstanceType) error {
	// Get regional instance types (default behavior - typically current region)
	regionalTypes, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{})
	if err != nil {
		return fmt.Errorf("failed to get regional instance types: %w", err)
	}

	if len(regionalTypes) == 0 {
		return errors.New("no regional instance types available for validation")
	}

	// Get all-region instance types by requesting from all locations
	allRegionTypes, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{
		Locations: All,
	})
	if err != nil {
		// If all-region is not supported, skip this validation
		return fmt.Errorf("all-region instance types not supported: %w", err)
	}

	if len(allRegionTypes) == 0 {
		return errors.New("no all-region instance types available for validation")
	}

	// Validate that regional results are a subset of all-region results
	if len(regionalTypes) >= len(allRegionTypes) {
		return fmt.Errorf("regional instance types (%d) should be fewer than all-region types (%d)",
			len(regionalTypes), len(allRegionTypes))
	}

	// Create a map of all-region types for efficient lookup
	allRegionMap := make(map[InstanceTypeID]InstanceType)
	for _, instanceType := range allRegionTypes {
		allRegionMap[instanceType.ID] = instanceType
	}

	// Validate that all regional types exist in all-region results
	for _, regionalType := range regionalTypes {
		if _, exists := allRegionMap[regionalType.ID]; !exists {
			return fmt.Errorf("regional instance type %s not found in all-region results", regionalType.ID)
		}
	}

	// Additional validation: ensure regional types have appropriate location information
	for _, regionalType := range regionalTypes {
		if regionalType.Location == "" {
			return fmt.Errorf("regional instance type %s should have location information", regionalType.ID)
		}
	}

	return nil
}

// normalizeInstanceTypes removes volatile fields that may change between calls
func normalizeInstanceTypes(types []InstanceType) []InstanceType {
	normalized := make([]InstanceType, len(types))
	for i, instanceType := range types {
		normalized[i] = instanceType
		// Remove fields that may vary between calls
		normalized[i].BasePrice = nil
		normalized[i].Quota = InstanceTypeQuota{}
		normalized[i].EstimatedDeployTime = nil
	}
	return normalized
}

// ValidateStableInstanceTypeIDs validates that the provided stable instance type IDs are valid and stable
// This function ensures that stable IDs exist in the current instance types and have required properties
func ValidateStableInstanceTypeIDs(ctx context.Context, client CloudInstanceType, stableIDs []InstanceTypeID) error {
	// Get all instance types
	allTypes, err := client.GetInstanceTypes(ctx, GetInstanceTypeArgs{})
	if err != nil {
		return fmt.Errorf("failed to get instance types: %w", err)
	}

	if len(allTypes) == 0 {
		return errors.New("no instance types available for validation")
	}

	// Group types by ID for efficient lookup
	typesByID := make(map[InstanceTypeID][]InstanceType)
	for _, instanceType := range allTypes {
		typesByID[instanceType.ID] = append(typesByID[instanceType.ID], instanceType)
	}

	// Validate that each ID has exactly one instance type (uniqueness)
	for id, types := range typesByID {
		if len(types) != 1 {
			return fmt.Errorf("instance type id %s should be unique, found %d instances", id, len(types))
		}
	}

	// Validate that stable IDs are not empty
	if len(stableIDs) == 0 {
		return errors.New("stable IDs list cannot be empty")
	}

	// Validate that all stable IDs exist in current instance types
	for _, stableID := range stableIDs {
		if _, exists := typesByID[stableID]; !exists {
			return fmt.Errorf("instance type id %s should be stable but not found", stableID)
		}
	}

	// Validate that all instance types have required properties
	for _, instanceType := range allTypes {
		// Check that instance type has base price
		if instanceType.BasePrice == nil {
			return fmt.Errorf("instance type %s should have base price", instanceType.ID)
		}

		// Check that supported storage has price information
		for i, storage := range instanceType.SupportedStorage {
			if storage.MinSize != nil {
				if storage.PricePerGBHr == nil {
					return fmt.Errorf("instance type %s should have storage %d price", instanceType.ID, i)
				}
			}
		}
	}

	return nil
}
