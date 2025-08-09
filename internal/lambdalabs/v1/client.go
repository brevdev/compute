package v1

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"github.com/cenkalti/backoff/v4"
)

// LambdaLabsCredential implements the CloudCredential interface for Lambda Labs
type LambdaLabsCredential struct {
	RefID  string `json:"ref_id"`
	APIKey string `json:"api_key"`
}

var _ v1.CloudCredential = &LambdaLabsCredential{}

// NewLambdaLabsCredential creates a new Lambda Labs credential
func NewLambdaLabsCredential(refID, apiKey string) *LambdaLabsCredential {
	return &LambdaLabsCredential{
		RefID:  refID,
		APIKey: apiKey,
	}
}

// GetReferenceID returns the reference ID for this credential
func (c *LambdaLabsCredential) GetReferenceID() string {
	return c.RefID
}

// GetAPIType returns the API type for Lambda Labs
func (c *LambdaLabsCredential) GetAPIType() v1.APIType {
	return v1.APITypeGlobal
}

const CloudProviderID = "lambda-labs"

const DefaultRegion string = "us-west-1"

// GetCloudProviderID returns the cloud provider ID for Lambda Labs
func (c *LambdaLabsCredential) GetCloudProviderID() v1.CloudProviderID {
	return CloudProviderID
}

// GetTenantID returns the tenant ID for Lambda Labs
func (c *LambdaLabsCredential) GetTenantID() (string, error) {
	return fmt.Sprintf("lambdalabs-%x", sha256.Sum256([]byte(c.APIKey))), nil
}

// MakeClient creates a new Lambda Labs client from this credential
func (c *LambdaLabsCredential) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
	return NewLambdaLabsClient(c.RefID, c.APIKey), nil
}

// LambdaLabsClient implements the CloudClient interface for Lambda Labs
// It embeds NotImplCloudClient to handle unsupported features
type LambdaLabsClient struct {
	v1.NotImplCloudClient
	refID    string
	apiKey   string
	baseURL  string
	client   *openapi.APIClient
	location string
}

var _ v1.CloudClient = &LambdaLabsClient{}

// NewLambdaLabsClient creates a new Lambda Labs client
func NewLambdaLabsClient(refID, apiKey string) *LambdaLabsClient {
	config := openapi.NewConfiguration()
	config.HTTPClient = http.DefaultClient
	client := openapi.NewAPIClient(config)

	return &LambdaLabsClient{
		refID:   refID,
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
func (c *LambdaLabsClient) MakeClient(_ context.Context, location string) (v1.CloudClient, error) {
	if location == "" {
		location = DefaultRegion
	}
	c.location = location
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
	return c.refID
}

func (c *LambdaLabsClient) makeAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, openapi.ContextBasicAuth, openapi.BasicAuth{
		UserName: c.apiKey,
	})
}

func getBackoff() backoff.BackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 1000 * time.Millisecond
	bo.MaxElapsedTime = 120 * time.Second
	return bo
}
