package v1

import (
	"context"

	"github.com/brevdev/cloud/pkg/v1"
)

func (c *FluidStackClient) CreateInstance(_ context.Context, _ v1.CreateInstanceAttrs) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) GetInstance(_ context.Context, _ v1.CloudProviderInstanceID) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) TerminateInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) ListInstances(_ context.Context, _ v1.ListInstancesArgs) ([]v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *FluidStackClient) StartInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) StopInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) RebootInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) ChangeInstanceType(_ context.Context, _ v1.CloudProviderInstanceID, _ string) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) UpdateInstanceTags(_ context.Context, _ v1.UpdateInstanceTagsArgs) error {
	return v1.ErrNotImplemented
}

func (c *FluidStackClient) MergeInstanceForUpdate(currInst v1.Instance, _ v1.Instance) v1.Instance {
	return currInst
}

func (c *FluidStackClient) MergeInstanceTypeForUpdate(currIt v1.InstanceType, _ v1.InstanceType) v1.InstanceType {
	return currIt
}
