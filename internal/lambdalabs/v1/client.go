package v1

import (
	"context"
	"net/http"

	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/compute/pkg/v1"
)

// LambdaLabsClient implements the CloudClient interface for Lambda Labs
// It embeds NotImplCloudClient to handle unsupported features
type LambdaLabsClient struct {
	v1.NotImplCloudClient
	apiKey  string
	baseURL string
	client  *openapi.APIClient
}

var _ v1.CloudClient = &LambdaLabsClient{}

// NewLambdaLabsClient creates a new Lambda Labs client
func NewLambdaLabsClient(apiKey string) *LambdaLabsClient {
	config := openapi.NewConfiguration()
	config.HTTPClient = http.DefaultClient
	client := openapi.NewAPIClient(config)

	return &LambdaLabsClient{
		apiKey:  apiKey,
		baseURL: "https://cloud.lambda.ai/api/v1",
		client:  client,
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
func (c *LambdaLabsClient) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
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

func (c *LambdaLabsClient) makeAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, openapi.ContextBasicAuth, openapi.BasicAuth{
		UserName: c.apiKey,
	})
}
