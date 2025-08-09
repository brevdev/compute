package commands

import (
	"context"
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/providers"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var createCmd = &cobra.Command{
	Use:   "create <refId>",
	Short: "Create a new instance",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreate,
}

var (
	instanceType string
	location     string
	name         string
	imageID      string
	publicKey    string
)

func init() {
	createCmd.Flags().StringVar(&instanceType, "instance-type", "", "Instance type to create")
	createCmd.Flags().StringVar(&location, "location", "", "Location to create instance in")
	createCmd.Flags().StringVar(&name, "name", "", "Name for the instance")
	createCmd.Flags().StringVar(&imageID, "image-id", "", "Image ID to use")
	createCmd.Flags().StringVar(&publicKey, "public-key", "", "SSH public key")

	if err := createCmd.MarkFlagRequired("instance-type"); err != nil {
		panic(err)
	}
}

func runCreate(_ *cobra.Command, args []string) error {
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	refID := args[0]

	credEntry, exists := cfg.Credentials[refID]
	if !exists {
		return fmt.Errorf("credential '%s' not found in config", refID)
	}

	cred, err := providers.CreateCredential(refID, credEntry)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	if location == "" {
		location = providers.GetDefaultLocation(cred)
	}
	if location == "" {
		return fmt.Errorf("location is required (use --location or set default_location in config)")
	}

	ctx := context.Background()
	client, err := cred.MakeClient(ctx, location)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	attrs := v1.CreateInstanceAttrs{
		Location:     location,
		Name:         name,
		InstanceType: instanceType,
		ImageID:      imageID,
		PublicKey:    publicKey,
	}

	instance, err := client.CreateInstance(ctx, attrs)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	output, err := yaml.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	fmt.Print(string(output))
	return nil
}
