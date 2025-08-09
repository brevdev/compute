package commands

import (
	"context"
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/config"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var typesCmd = &cobra.Command{
	Use:   "types <refId>",
	Short: "List available instance types",
	Args:  cobra.ExactArgs(1),
	RunE:  runTypes,
}

var typesLocation string

func init() {
	typesCmd.Flags().StringVar(&typesLocation, "location", "", "Location to get instance types for")
}

func runTypes(_ *cobra.Command, args []string) error {
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	refID := args[0]

	credEntry, exists := cfg.Credentials[refID]
	if !exists {
		return fmt.Errorf("credential '%s' not found in config", refID)
	}

	cred := credEntry.Value
	if cred == nil {
		return fmt.Errorf("credential value is nil")
	}

	if typesLocation == "" {
		if provider, ok := cred.(config.DefaultLocationProvider); ok {
			typesLocation = provider.GetDefaultLocation()
		}
	}
	if typesLocation == "" {
		return fmt.Errorf("location is required (use --location or set default_location in config)")
	}

	ctx := context.Background()
	client, err := cred.MakeClient(ctx, typesLocation)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	instanceTypes, err := client.GetInstanceTypes(ctx, v1.GetInstanceTypeArgs{})
	if err != nil {
		return fmt.Errorf("failed to get instance types: %w", err)
	}

	output, err := yaml.Marshal(instanceTypes)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Print(string(output))
	return nil
}
