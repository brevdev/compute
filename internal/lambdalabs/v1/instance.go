package v1

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// CreateInstance creates a new instance in Lambda Labs
// Supported via: POST /api/v1/instance-operations/launch
func (c *LambdaLabsClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	// TODO: Implement Lambda Labs instance creation
	// This would typically involve:
	// 1. Validating the instance type and location
	// 2. Creating the instance via Lambda Labs API
	// 3. Waiting for the instance to be ready
	// 4. Returning the instance details

	return &v1.Instance{
		Name:         attrs.Name,
		RefID:        attrs.RefID,
		CreatedAt:    time.Now(),
		CloudID:      v1.CloudProviderInstanceID("lambda-instance-id"), // TODO: Get from API response
		Location:     attrs.Location,
		SubLocation:  attrs.SubLocation,
		InstanceType: attrs.InstanceType,
		ImageID:      attrs.ImageID,
		DiskSize:     attrs.DiskSize,
		Status: v1.Status{
			LifecycleStatus: v1.LifecycleStatusRunning,
		},
		Tags: attrs.Tags,
	}, nil
}

// GetInstance retrieves an instance by ID
// Supported via: GET /api/v1/instances/{id}
func (c *LambdaLabsClient) GetInstance(ctx context.Context, id v1.CloudProviderInstanceID) (*v1.Instance, error) {
	// TODO: Implement Lambda Labs instance retrieval
	return nil, fmt.Errorf("not implemented")
}

// TerminateInstance terminates an instance
// Supported via: POST /api/v1/instance-operations/terminate
func (c *LambdaLabsClient) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	// TODO: Implement Lambda Labs instance termination
	return fmt.Errorf("not implemented")
}

// ListInstances lists all instances
// Supported via: GET /api/v1/instances
func (c *LambdaLabsClient) ListInstances(ctx context.Context, args v1.ListInstancesArgs) ([]v1.Instance, error) {
	// TODO: Implement Lambda Labs instance listing
	return nil, fmt.Errorf("not implemented")
}

// RebootInstance reboots an instance
// Supported via: POST /api/v1/instance-operations/restart
func (c *LambdaLabsClient) RebootInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	// TODO: Implement Lambda Labs instance rebooting
	return fmt.Errorf("not implemented")
}

// UNSUPPORTED OPERATIONS (no API endpoints found):

// StopInstance stops an instance
// UNSUPPORTED: No stop endpoint found in Lambda Labs API
func (c *LambdaLabsClient) StopInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

// StartInstance starts an instance
// UNSUPPORTED: No start endpoint found in Lambda Labs API
func (c *LambdaLabsClient) StartInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

// ChangeInstanceType changes the instance type
// UNSUPPORTED: No instance type modification endpoint found in Lambda Labs API
func (c *LambdaLabsClient) ChangeInstanceType(ctx context.Context, instanceID v1.CloudProviderInstanceID, instanceType string) error {
	return v1.ErrNotImplemented
}

// UpdateInstanceTags updates instance tags
// UNSUPPORTED: No tagging endpoints found in Lambda Labs API
func (c *LambdaLabsClient) UpdateInstanceTags(ctx context.Context, args v1.UpdateInstanceTagsArgs) error {
	return v1.ErrNotImplemented
}

// MergeInstanceForUpdate merges instance data for updates
func (c *LambdaLabsClient) MergeInstanceForUpdate(currInst v1.Instance, newInst v1.Instance) v1.Instance {
	// TODO: Implement instance merging logic
	return newInst
}

// MergeInstanceTypeForUpdate merges instance type data for updates
func (c *LambdaLabsClient) MergeInstanceTypeForUpdate(currIt v1.InstanceType, newIt v1.InstanceType) v1.InstanceType {
	// TODO: Implement instance type merging logic
	return newIt
}
