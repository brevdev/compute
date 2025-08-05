package v1

import (
	"context"
	"net/http"

	"github.com/brevdev/cloud/pkg/v1"
)

type FluidStackClient struct {
	*v1.NotImplCloudClient
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewFluidStackClient(apiKey string) *FluidStackClient {
	return &FluidStackClient{
		NotImplCloudClient: &v1.NotImplCloudClient{},
		baseURL:            "https://api.fluidstack.io/v1alpha1",
		httpClient:         &http.Client{},
		apiKey:             apiKey,
	}
}

func (c *FluidStackClient) GetCapabilities(_ context.Context) (v1.Capabilities, error) {
	return []v1.Capability{
		v1.CapabilityCreateInstance,
		v1.CapabilityTerminateInstance,
		v1.CapabilityStopStartInstance,
		v1.CapabilityTags,
		CapabilityCreateProject,
		CapabilityDeleteProject,
		CapabilityListProjects,
		CapabilityGetProject,
	}, nil
}

func (c *FluidStackClient) GetName() string {
	return "fluidstack"
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

func (c *FluidStackClient) GetRegions(_ context.Context) ([]*v1.Location, error) {
	return nil, v1.ErrNotImplemented
}
