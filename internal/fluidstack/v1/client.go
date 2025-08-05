package v1

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/brevdev/cloud/pkg/v1"
)

// FluidStackCredential implements the CloudCredential interface for FluidStack
type FluidStackCredential struct {
	RefID  string
	APIKey string
}

var _ v1.CloudCredential = &FluidStackCredential{}

func NewFluidStackCredential(refID, apiKey string) *FluidStackCredential {
	return &FluidStackCredential{
		RefID:  refID,
		APIKey: apiKey,
	}
}

// GetReferenceID returns the reference ID for this credential
func (c *FluidStackCredential) GetReferenceID() string {
	return c.RefID
}

// GetAPIType returns the API type for FluidStack
func (c *FluidStackCredential) GetAPIType() v1.APIType {
	return v1.APITypeGlobal
}

// GetCloudProviderID returns the cloud provider ID for FluidStack
func (c *FluidStackCredential) GetCloudProviderID() v1.CloudProviderID {
	return "fluidstack"
}

// GetTenantID returns the tenant ID for FluidStack
func (c *FluidStackCredential) GetTenantID() (string, error) {
	return fmt.Sprintf("fluidstack-%x", sha256.Sum256([]byte(c.APIKey))), nil
}

// GetCapabilities returns the capabilities for FluidStack
func (c *FluidStackCredential) GetCapabilities(ctx context.Context) (v1.Capabilities, error) {
	client, err := c.MakeClient(ctx, "")
	if err != nil {
		return nil, err
	}
	return client.GetCapabilities(ctx)
}

// MakeClient creates a new FluidStack client from this credential
func (c *FluidStackCredential) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
	return NewFluidStackClient(c.RefID, c.APIKey), nil
}

// FluidStackClient implements the CloudClient interface for FluidStack
// It embeds NotImplCloudClient to handle unsupported features
type FluidStackClient struct {
	v1.NotImplCloudClient
	refID   string
	apiKey  string
	baseURL string
}

var _ v1.CloudClient = &FluidStackClient{}

func NewFluidStackClient(refID, apiKey string) *FluidStackClient {
	return &FluidStackClient{
		refID:   refID,
		apiKey:  apiKey,
		baseURL: "https://api.fluidstack.io/v1alpha1",
	}
}

// GetAPIType returns the API type for FluidStack
func (c *FluidStackClient) GetAPIType() v1.APIType {
	return v1.APITypeGlobal
}

// GetCloudProviderID returns the cloud provider ID for FluidStack
func (c *FluidStackClient) GetCloudProviderID() v1.CloudProviderID {
	return "fluidstack"
}

// MakeClient creates a new client instance
func (c *FluidStackClient) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
	return c, nil
}

// GetReferenceID returns the reference ID for this client
func (c *FluidStackClient) GetReferenceID() string {
	return c.refID
}
