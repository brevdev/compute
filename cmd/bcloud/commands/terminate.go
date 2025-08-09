package commands

import (
	"context"
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/providers"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/spf13/cobra"
)

var terminateCmd = &cobra.Command{
	Use:   "terminate <refId> <instanceId>",
	Short: "Terminate an instance",
	Args:  cobra.ExactArgs(2),
	RunE:  runTerminate,
}

func runTerminate(_ *cobra.Command, args []string) error {
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	refID := args[0]
	instanceID := v1.CloudProviderInstanceID(args[1])

	credEntry, exists := cfg.Credentials[refID]
	if !exists {
		return fmt.Errorf("credential '%s' not found in config", refID)
	}

	cred, err := providers.CreateCredential(refID, credEntry)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	defaultLocation := providers.GetDefaultLocation(cred)
	if defaultLocation == "" {
		return fmt.Errorf("default location is required in config")
	}

	ctx := context.Background()
	client, err := cred.MakeClient(ctx, defaultLocation)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	err = client.TerminateInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to terminate instance: %w", err)
	}

	fmt.Printf("Instance %s terminated successfully\n", instanceID)
	return nil
}
