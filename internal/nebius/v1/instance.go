package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *NebiusClient) GetInstance(ctx context.Context, id v1.CloudProviderInstanceID) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *NebiusClient) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) ListInstances(ctx context.Context, args v1.ListInstancesArgs) ([]v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *NebiusClient) StopInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) StartInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) RebootInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) MergeInstanceForUpdate(currInst v1.Instance, newInst v1.Instance) v1.Instance {
	merged := newInst
	
	merged.Name = currInst.Name
	merged.RefID = currInst.RefID
	merged.CloudCredRefID = currInst.CloudCredRefID
	merged.CreatedAt = currInst.CreatedAt
	merged.CloudID = currInst.CloudID
	merged.Location = currInst.Location
	merged.SubLocation = currInst.SubLocation
	
	return merged
}
