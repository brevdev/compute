package v1

import (
	"context"
	"fmt"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func (c *NebiusClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	securityGroupID, err := c.ensureClusterSecurityGroup(ctx, attrs)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure cluster security group: %w", err)
	}

	instance, err := c.createInstanceWithSecurityGroup(ctx, attrs, securityGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance with security group: %w", err)
	}

	return instance, nil
}

func (c *NebiusClient) ensureClusterSecurityGroup(_ context.Context, attrs v1.CreateInstanceAttrs) (string, error) {
	clusterID := c.getClusterIDFromAttrs(attrs)
	_ = fmt.Sprintf("brev-cluster-%s", clusterID)

	return "", fmt.Errorf("cluster security group creation not yet implemented - need to use Nebius VPC service")
}

func (c *NebiusClient) createInstanceWithSecurityGroup(_ context.Context, _ v1.CreateInstanceAttrs, _ string) (*v1.Instance, error) {
	return nil, fmt.Errorf("instance creation with security group not yet implemented - need to use Nebius Compute service")
}

func (c *NebiusClient) getClusterIDFromAttrs(attrs v1.CreateInstanceAttrs) string {
	if attrs.Tags != nil {
		if clusterID, exists := attrs.Tags["cluster_id"]; exists {
			return clusterID
		}
	}
	return "default"
}

func (c *NebiusClient) GetInstance(_ context.Context, _ v1.CloudProviderInstanceID) (*v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *NebiusClient) TerminateInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) ListInstances(_ context.Context, _ v1.ListInstancesArgs) ([]v1.Instance, error) {
	return nil, v1.ErrNotImplemented
}

func (c *NebiusClient) StopInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) StartInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *NebiusClient) RebootInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
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
