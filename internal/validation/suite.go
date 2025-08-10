package validation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brevdev/cloud/pkg/ssh"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/stretchr/testify/require"
)

const UniversalCloudCredRefID = "brev-validation-test"

type ProviderConfig struct {
	Location   string
	StableIDs  []v1.InstanceTypeID
	Credential v1.CloudCredential
}

func RunValidationSuite(t *testing.T, config ProviderConfig) {
	if testing.Short() {
		t.Skip("Skipping validation tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client, err := config.Credential.MakeClient(ctx, config.Location)
	if err != nil {
		t.Fatalf("Failed to create client for %s: %v", config.Credential.GetCloudProviderID(), err)
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
		err := v1.ValidateLocationalInstanceTypes(ctx, client)
		require.NoError(t, err, "ValidateRegionalInstanceTypes should pass")
	})

	t.Run("ValidateStableInstanceTypeIDs", func(t *testing.T) {
		err = v1.ValidateStableInstanceTypeIDs(ctx, client, config.StableIDs)
		require.NoError(t, err, "ValidateStableInstanceTypeIDs should pass")
	})
}

func RunInstanceLifecycleValidation(t *testing.T, config ProviderConfig) {
	if testing.Short() {
		t.Skip("Skipping validation tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client, err := config.Credential.MakeClient(ctx, config.Location)
	if err != nil {
		t.Fatalf("Failed to create client for %s: %v", config.Credential.GetCloudProviderID(), err)
	}
	capabilities, err := client.GetCapabilities(ctx)
	require.NoError(t, err)

	types, err := client.GetInstanceTypes(ctx, v1.GetInstanceTypeArgs{})
	require.NoError(t, err)
	require.NotEmpty(t, types, "Should have instance types")

	locations, err := client.GetLocations(ctx, v1.GetLocationsArgs{})
	require.NoError(t, err)
	require.NotEmpty(t, locations, "Should have locations")

	t.Run("ValidateCreateInstance", func(t *testing.T) {
		attrs := v1.CreateInstanceAttrs{}
		for _, typ := range types {
			if typ.IsAvailable {
				attrs.InstanceType = typ.Type
				attrs.Location = typ.Location
				attrs.PublicKey = ssh.GetTestPublicKey()
				break
			}
		}
		instance, err := v1.ValidateCreateInstance(ctx, client, attrs)
		if err != nil {
			t.Fatalf("ValidateCreateInstance failed: %v", err)
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

		t.Run("ValidateSSHAccessible", func(t *testing.T) {
			err := v1.ValidateInstanceSSHAccessible(ctx, client, instance, ssh.GetTestPrivateKey())
			require.NoError(t, err, "ValidateSSHAccessible should pass")
		})

		instance, err = client.GetInstance(ctx, instance.CloudID)
		require.NoError(t, err)

		t.Run("ValidateInstanceImage", func(t *testing.T) {
			err := v1.ValidateInstanceImage(ctx, *instance, ssh.GetTestPrivateKey())
			require.NoError(t, err, "ValidateInstanceImage should pass")
		})

		if capabilities.IsCapable(v1.CapabilityStopStartInstance) && instance.Stoppable {
			t.Run("ValidateStopStartInstance", func(t *testing.T) {
				err := v1.ValidateStopStartInstance(ctx, client, instance)
				require.NoError(t, err, "ValidateStopStartInstance should pass")
			})
		}

		t.Run("ValidateTerminateInstance", func(t *testing.T) {
			err := v1.ValidateTerminateInstance(ctx, client, instance)
			require.NoError(t, err, "ValidateTerminateInstance should pass")
		})
	})
}

func CleanupOrphanedInstances(ctx context.Context, client v1.CloudCreateTerminateInstance) error {
	instances, err := client.ListInstances(ctx, v1.ListInstancesArgs{})
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	cutoffTime := time.Now().Add(-1 * time.Hour)
	var orphanedInstances []v1.Instance

	for _, instance := range instances {
		if instance.CloudCredRefID == UniversalCloudCredRefID {
			if instance.CreatedAt.Before(cutoffTime) {
				orphanedInstances = append(orphanedInstances, instance)
			}
		}
	}

	if len(orphanedInstances) == 0 {
		fmt.Printf("No orphaned instances found with CloudCredRefID: %s\n", UniversalCloudCredRefID)
		return nil
	}

	fmt.Printf("Found %d orphaned instances to clean up\n", len(orphanedInstances))

	var cleanupErrors []error
	for _, instance := range orphanedInstances {
		fmt.Printf("Terminating orphaned instance: %s (created: %s)\n",
			instance.CloudID, instance.CreatedAt.Format(time.RFC3339))

		err := client.TerminateInstance(ctx, instance.CloudID)
		if err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Errorf("failed to terminate instance %s: %w", instance.CloudID, err))
		}
	}

	if len(cleanupErrors) > 0 {
		return fmt.Errorf("cleanup completed with %d errors: %v", len(cleanupErrors), cleanupErrors)
	}

	fmt.Printf("Successfully cleaned up %d orphaned instances\n", len(orphanedInstances))
	return nil
}
