package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// LambdaLabsClient implements the CloudClient interface for Lambda Labs
type LambdaLabsClient struct {
	// TODO: Add Lambda Labs specific fields like API key, base URL, etc.
	apiKey  string
	baseURL string
}

var _ v1.CloudClient = &LambdaLabsClient{}

// NewLambdaLabsClient creates a new Lambda Labs client
func NewLambdaLabsClient(apiKey string) *LambdaLabsClient {
	return &LambdaLabsClient{
		apiKey:  apiKey,
		baseURL: "https://cloud.lambdalabs.com/api/v1",
	}
}

// GetAPIType returns the API type for Lambda Labs
func (c *LambdaLabsClient) GetAPIType() v1.APIType {
	return v1.APITypeGlobal // Lambda Labs uses a global API
}

// GetCloudProviderID returns the cloud provider ID for Lambda Labs
func (c *LambdaLabsClient) GetCloudProviderID() v1.CloudProviderID {
	return "lambdalabs"
}

// MakeClient creates a new client instance
func (c *LambdaLabsClient) MakeClient(ctx context.Context, location string) (v1.CloudClient, error) {
	// Lambda Labs doesn't require location-specific clients
	return c, nil
}

// GetTenantID returns the tenant ID for Lambda Labs
func (c *LambdaLabsClient) GetTenantID() (string, error) {
	// TODO: Implement tenant ID retrieval for Lambda Labs
	return "", nil
}

// GetReferenceID returns the reference ID for this client
func (c *LambdaLabsClient) GetReferenceID() string {
	// TODO: Implement reference ID generation
	return "lambdalabs-client"
}
