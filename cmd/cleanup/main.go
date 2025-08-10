package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	lambdalabs "github.com/brevdev/cloud/internal/lambdalabs/v1"
	"github.com/brevdev/cloud/internal/validation"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

func main() {
	var (
		provider = flag.String("provider", "", "Cloud provider to clean up (lambdalabs)")
		dryRun   = flag.Bool("dry-run", false, "List orphaned instances without deleting them")
	)
	flag.Parse()

	if *provider == "" {
		log.Fatal("Provider is required. Use -provider=lambdalabs")
	}

	if *provider != "lambdalabs" {
		log.Fatalf("Unsupported provider: %s", *provider)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := cleanupLambdaLabs(ctx, *dryRun)
	if err != nil {
		log.Printf("LambdaLabs cleanup failed: %v", err)
		return
	}
}

func cleanupLambdaLabs(ctx context.Context, dryRun bool) error {
	apiKey := os.Getenv("LAMBDALABS_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("LAMBDALABS_API_KEY environment variable is required")
	}

	credential := lambdalabs.NewLambdaLabsCredential(validation.UniversalCloudCredRefID, apiKey)
	client, err := credential.MakeClient(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to create LambdaLabs client: %w", err)
	}

	if dryRun {
		return listOrphanedInstances(ctx, client)
	}

	return validation.CleanupOrphanedInstances(ctx, client)
}

func listOrphanedInstances(ctx context.Context, client v1.CloudClient) error {
	instances, err := client.ListInstances(ctx, v1.ListInstancesArgs{})
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	cutoffTime := time.Now().Add(-1 * time.Hour)
	var orphanedInstances []v1.Instance

	for _, instance := range instances {
		if instance.CloudCredRefID == validation.UniversalCloudCredRefID {
			if instance.CreatedAt.Before(cutoffTime) {
				orphanedInstances = append(orphanedInstances, instance)
			}
		}
	}

	fmt.Printf("Found %d orphaned instances with CloudCredRefID: %s\n",
		len(orphanedInstances), validation.UniversalCloudCredRefID)

	for _, instance := range orphanedInstances {
		fmt.Printf("- Instance: %s (created: %s, age: %s)\n",
			instance.CloudID,
			instance.CreatedAt.Format(time.RFC3339),
			time.Since(instance.CreatedAt).Round(time.Minute))
	}

	return nil
}
