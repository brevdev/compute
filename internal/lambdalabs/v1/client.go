package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// LambdaLabsClient implements the CloudClient interface for Lambda Labs
// It embeds NotImplCloudClient to handle unsupported features
type LambdaLabsClient struct {
	v1.NotImplCloudClient
	apiKey  string
	baseURL string
}

var _ v1.CloudClient = &LambdaLabsClient{}

// NewLambdaLabsClient creates a new Lambda Labs client
func NewLambdaLabsClient(apiKey string) *LambdaLabsClient {
	return &LambdaLabsClient{
		apiKey:  apiKey,
		baseURL: "https://cloud.lambda.ai/api/v1",
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
	// This could be derived from the API key or account information
	return "", nil
}

// GetReferenceID returns the reference ID for this client
func (c *LambdaLabsClient) GetReferenceID() string {
	return "lambdalabs-client"
}
