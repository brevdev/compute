package commands

import (
	"context"
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/config"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var listCmd = &cobra.Command{
	Use:   "list <refId>",
	Short: "List instances",
	Args:  cobra.ExactArgs(1),
	RunE:  runList,
}

var listLocation string

func init() {
	listCmd.Flags().StringVar(&listLocation, "location", "", "Location to list instances from")
}

func runList(_ *cobra.Command, args []string) error {
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
		return fmt.Errorf("credential entry has no value")
	}

	if listLocation == "" {
		if provider, ok := cred.(config.DefaultLocationProvider); ok {
			listLocation = provider.GetDefaultLocation()
		}
	}
	if listLocation == "" {
		return fmt.Errorf("location is required (use --location or set default_location in config)")
	}

	ctx := context.Background()
	client, err := cred.MakeClient(ctx, listLocation)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	instances, err := client.ListInstances(ctx, v1.ListInstancesArgs{
		Locations: []string{listLocation},
	})
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	output, err := yaml.Marshal(instances)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Print(string(output))
	return nil
}
