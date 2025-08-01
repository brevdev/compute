package v1

import (
	"context"
	"time"
)

type APIType string

const (
	APITypeLocational APIType = "locational"
	APITypeGlobal     APIType = "global"
)

type CloudProviderID string

type CloudUtils interface {
	GetCapabilities(ctx context.Context) (Capabilities, error)
	GetCloudProviderID() CloudProviderID
}

type CloudCredential interface {
	GetAPIType() APIType
	MakeClient(ctx context.Context, location string) (CloudClient, error)
	GetTenantID() (string, error)
	GetReferenceID() string
	CloudUtils
	GetInstanceTypePollTime() time.Duration
}

type CloudBase interface {
	CloudCreateTerminateInstance
	CloudLocation
}

type CloudClient interface {
	CloudCredential
	CloudBase
	CloudQuota
	CloudRebootInstance
	CloudStopStartInstance
	CloudResizeInstanceVolume
	CloudMachineImage
	CloudBillingUsage
	CloudChangeInstanceType
	CloudRetireVolume
	CloudVPCSubnets
	CloudModifyFirewall
	CloudSnapshotter
	CloudCreateMachineImage
	CloudInstanceTags
	UpdateHandler
}
