package config

import (
	"context"
	"fmt"

	v1 "github.com/brevdev/cloud/pkg/v1"

	lambdalabs "github.com/brevdev/cloud/internal/lambdalabs/v1"
)

type LambdaLabsCredential struct {
	APIKey          string `json:"api_key" yaml:"api_key"`
	DefaultLocation string `json:"default_location" yaml:"default_location"`
	RefID           string `json:"ref_id,omitempty" yaml:"ref_id,omitempty"`
}

func (l *LambdaLabsCredential) GetCloudProviderID() v1.CloudProviderID {
	return "lambdalabs"
}

func (l *LambdaLabsCredential) GetReferenceID() string {
	return l.RefID
}

func (l *LambdaLabsCredential) GetAPIType() v1.APIType {
	return v1.APITypeLocational
}

func (l *LambdaLabsCredential) GetTenantID() (string, error) {
	return "", nil
}

func (l *LambdaLabsCredential) GetCapabilities(ctx context.Context) (v1.Capabilities, error) {
	cred := lambdalabs.NewLambdaLabsCredential(l.RefID, l.APIKey)
	return cred.GetCapabilities(ctx)
}

func (l *LambdaLabsCredential) MakeClient(ctx context.Context, location string) (v1.CloudClient, error) {
	actualRefID := l.RefID
	if actualRefID == "" {
		return nil, fmt.Errorf("ref_id is required")
	}
	cred := lambdalabs.NewLambdaLabsCredential(actualRefID, l.APIKey)
	return cred.MakeClient(ctx, location)
}

func (l *LambdaLabsCredential) GetDefaultLocation() string {
	return l.DefaultLocation
}

type NebiusCredential struct {
	ServiceAccountJSON string `json:"service_account_json" yaml:"service_account_json"`
	DefaultLocation    string `json:"default_location" yaml:"default_location"`
	RefID              string `json:"ref_id,omitempty" yaml:"ref_id,omitempty"`
}

func (n *NebiusCredential) GetCloudProviderID() v1.CloudProviderID {
	return "nebius"
}

func (n *NebiusCredential) GetReferenceID() string {
	return n.RefID
}

func (n *NebiusCredential) GetAPIType() v1.APIType {
	return v1.APITypeLocational
}

func (n *NebiusCredential) GetTenantID() (string, error) {
	return "", nil
}

func (n *NebiusCredential) GetCapabilities(_ context.Context) (v1.Capabilities, error) {
	return v1.Capabilities{}, nil
}

func (n *NebiusCredential) MakeClient(_ context.Context, _ string) (v1.CloudClient, error) {
	return nil, fmt.Errorf("nebius provider not yet implemented")
}

func (n *NebiusCredential) GetDefaultLocation() string {
	return n.DefaultLocation
}

func init() {
	RegisterProvider("lambdalabs", func() v1.CloudCredential { return &LambdaLabsCredential{} })
	RegisterProvider("nebius", func() v1.CloudCredential { return &NebiusCredential{} })
}
