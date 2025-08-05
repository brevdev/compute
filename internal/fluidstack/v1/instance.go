package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) GetInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) ListInstances(ctx context.Context, args v1.ListInstancesArgs) ([]v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) StartInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) StopInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) RebootInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) ChangeInstanceType(ctx context.Context, instanceID v1.CloudProviderInstanceID, instanceType string) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) UpdateInstanceTags(ctx context.Context, args v1.UpdateInstanceTagsArgs) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) MergeInstanceForUpdate(currInst v1.Instance, newInst v1.Instance) v1.Instance {
	return currInst
}

func (c *FluidStackClient) MergeInstanceTypeForUpdate(currIt v1.InstanceType, newIt v1.InstanceType) v1.InstanceType {
	return currIt
}
