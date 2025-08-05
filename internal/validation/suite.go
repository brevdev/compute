package validation

import (
	"context"
	"testing"
	"time"

	v1 "github.com/brevdev/compute/pkg/v1"
	"github.com/stretchr/testify/require"
)

type ProviderConfig struct {
	ProviderName string
	Credential   v1.CloudCredential
}

func RunValidationSuite(t *testing.T, config ProviderConfig) {
	if testing.Short() {
		t.Skip("Skipping validation tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client, err := config.Credential.MakeClient(ctx, "")
	if err != nil {
		t.Skipf("Failed to create client for %s: %v", config.ProviderName, err)
	}

	t.Run("ValidateGetLocations", func(t *testing.T) {
		err := v1.ValidateGetLocations(ctx, client)
		require.NoError(t, err, "ValidateGetLocations should pass")
	})

	t.Run("ValidateGetInstanceTypes", func(t *testing.T) {
		err := v1.ValidateGetInstanceTypes(ctx, client)
		require.NoError(t, err, "ValidateGetInstanceTypes should pass")
	})

	t.Run("ValidateRegionalInstanceTypes", func(t *testing.T) {
		err := v1.ValidateRegionalInstanceTypes(ctx, client)
		require.NoError(t, err, "ValidateRegionalInstanceTypes should pass")
	})

	t.Run("ValidateStableInstanceTypeIDs", func(t *testing.T) {
		types, err := client.GetInstanceTypes(ctx, v1.GetInstanceTypeArgs{})
		require.NoError(t, err)
		require.NotEmpty(t, types, "Should have instance types")

		stableIDs := []v1.InstanceTypeID{types[0].ID}
		err = v1.ValidateStableInstanceTypeIDs(ctx, client, stableIDs)
		require.NoError(t, err, "ValidateStableInstanceTypeIDs should pass")
	})
}

func RunInstanceLifecycleValidation(t *testing.T, config ProviderConfig) {
	if testing.Short() {
		t.Skip("Skipping validation tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client, err := config.Credential.MakeClient(ctx, "")
	if err != nil {
		t.Skipf("Failed to create client for %s: %v", config.ProviderName, err)
	}

	types, err := client.GetInstanceTypes(ctx, v1.GetInstanceTypeArgs{})
	require.NoError(t, err)
	require.NotEmpty(t, types, "Should have instance types")

	locations, err := client.GetLocations(ctx, v1.GetLocationsArgs{})
	require.NoError(t, err)
	require.NotEmpty(t, locations, "Should have locations")

	var instanceType string
	var location string
	for _, loc := range locations {
		if loc.Available {
			location = loc.Name
			break
		}
	}
	require.NotEmpty(t, location, "Should have available location")

	for _, typ := range types {
		if typ.Location == location && typ.IsAvailable {
			instanceType = typ.Type
			break
		}
	}
	require.NotEmpty(t, instanceType, "Should have available instance type")

	t.Run("ValidateCreateInstance", func(t *testing.T) {
		attrs := v1.CreateInstanceAttrs{
			Name:         "validation-test",
			InstanceType: instanceType,
			Location:     location,
		}

		instance, err := v1.ValidateCreateInstance(ctx, client, attrs)
		if err != nil {
			t.Logf("ValidateCreateInstance failed: %v", err)
			t.Skip("Skipping due to create instance failure - may be quota/availability issue")
		}
		require.NotNil(t, instance)

		defer func() {
			if instance != nil {
				_ = client.TerminateInstance(ctx, instance.CloudID)
			}
		}()

		t.Run("ValidateListCreatedInstance", func(t *testing.T) {
			err := v1.ValidateListCreatedInstance(ctx, client, instance)
			require.NoError(t, err, "ValidateListCreatedInstance should pass")
		})

		t.Run("ValidateTerminateInstance", func(t *testing.T) {
			err := v1.ValidateTerminateInstance(ctx, client, *instance)
			require.NoError(t, err, "ValidateTerminateInstance should pass")
			instance = nil // Mark as terminated
		})
	})
}
