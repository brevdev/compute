package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

type FluidStackClient struct {
	v1.NotImplCloudClient
	apiKey  string
	baseURL string
}

var _ v1.CloudClient = &FluidStackClient{}

func NewFluidStackClient(apiKey string) *FluidStackClient {
	return &FluidStackClient{
		apiKey:  apiKey,
		baseURL: "https://api.fluidstack.io/v1alpha1",
	}
}

func (c *FluidStackClient) GetAPIType() v1.APIType {
	return v1.APITypeGlobal
}

func (c *FluidStackClient) GetCloudProviderID() v1.CloudProviderID {
	return "fluidstack"
}

func (c *FluidStackClient) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
	return c, nil
}

func (c *FluidStackClient) GetTenantID() (string, error) {
	return "", v1.ErrNotImplemented
}

func (c *FluidStackClient) GetReferenceID() string {
	return c.apiKey
}
