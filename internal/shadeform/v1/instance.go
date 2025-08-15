package v1

import (
	"context"
	"fmt"
	"github.com/alecthomas/units"
	openapi "github.com/brevdev/cloud/internal/shadeform/gen/shadeform"
	"github.com/brevdev/cloud/pkg/v1"
)

const (
	hostname = "shadecloud"
)

func (c *ShadeformClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	authCtx := c.makeAuthContext(ctx)

	region := attrs.Location
	cloud, shadeInstanceType, err := c.getShadeformCloudAndInstanceType(attrs.InstanceType)
	if err != nil {
		return nil, err
	}

	cloudEnum, err := openapi.NewCloudFromValue(cloud)
	if err != nil {
		return nil, err
	}

	req := openapi.CreateRequest{
		Cloud:             *cloudEnum,
		Region:            region,
		ShadeInstanceType: shadeInstanceType,
		Name:              attrs.Name,
	}

	resp, httpResp, err := c.client.DefaultAPI.InstancesCreate(authCtx).CreateRequest(req).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("no instance returned from create request")
	}

	// Since Shadeform doesn't return the full instance that's created, we need to make a second API call to get
	// the created instance after we call create
	createdInstance, err := c.GetInstance(authCtx, v1.CloudProviderInstanceID(resp.Id))
	if err != nil {
		return nil, err
	}

	return createdInstance, nil
}

func (c *ShadeformClient) GetInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) (*v1.Instance, error) {
	authCtx := c.makeAuthContext(ctx)

	resp, httpResp, err := c.client.DefaultAPI.InstancesInfo(authCtx, string(instanceID)).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("no instance returned from get request")
	}

	return c.convertInstanceInfoResponseToV1Instance(*resp), nil
}

func (c *ShadeformClient) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	authCtx := c.makeAuthContext(ctx)

	httpResp, err := c.client.DefaultAPI.InstancesDelete(authCtx, string(instanceID)).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return fmt.Errorf("failed to terminate instance: %w", err)
	}

	return nil
}

func (c *ShadeformClient) ListInstances(ctx context.Context, _ v1.ListInstancesArgs) ([]v1.Instance, error) {
	authCtx := c.makeAuthContext(ctx)

	resp, httpResp, err := c.client.DefaultAPI.Instances(authCtx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var instances []v1.Instance
	for _, instance := range resp.Instances {
		instances = append(instances, *c.convertShadeformInstanceToV1Instance(instance))
	}

	return instances, nil
}

func (c *ShadeformClient) RebootInstance(_ context.Context, _ v1.CloudProviderInstanceID) error {
	return v1.ErrNotImplemented
}

func (c *ShadeformClient) MergeInstanceForUpdate(currInst v1.Instance, _ v1.Instance) v1.Instance {
	return currInst
}

func (c *ShadeformClient) MergeInstanceTypeForUpdate(currIt v1.InstanceType, _ v1.InstanceType) v1.InstanceType {
	return currIt
}

func (c *ShadeformClient) convertInstanceInfoResponseToV1Instance(instanceInfo openapi.InstanceInfoResponse) *v1.Instance {
	var lifecycleStatus v1.LifecycleStatus
	switch instanceInfo.Status {
	case "creating":
		lifecycleStatus = v1.LifecycleStatusPending
	case "pending_provider":
		lifecycleStatus = v1.LifecycleStatusPending
	case "pending":
		lifecycleStatus = v1.LifecycleStatusPending
	case "active":
		lifecycleStatus = v1.LifecycleStatusRunning
	case "error":
		lifecycleStatus = v1.LifecycleStatusFailed
	default:
		lifecycleStatus = v1.LifecycleStatusPending
	}

	instanceType := c.getInstanceType(string(instanceInfo.Cloud), instanceInfo.ShadeInstanceType)

	instance := &v1.Instance{
		Name:         instanceInfo.Name,
		CreatedAt:    instanceInfo.CreatedAt,
		CloudID:      v1.CloudProviderInstanceID(instanceInfo.Id),
		PublicIP:     instanceInfo.Ip,
		Hostname:     hostname,
		ImageID:      instanceInfo.Configuration.Os,
		InstanceType: instanceType,
		DiskSize:     units.Base2Bytes(instanceInfo.Configuration.StorageInGb) * units.GiB,
		SSHUser:      instanceInfo.SshUser,
		SSHPort:      int(instanceInfo.SshPort),
		Status: v1.Status{
			LifecycleStatus: lifecycleStatus,
		},
		Spot:       false,
		Location:   instanceInfo.Region,
		Stoppable:  false,
		Rebootable: true,
	}

	return instance
}

func (c *ShadeformClient) convertShadeformInstanceToV1Instance(shadeformInstance openapi.Instance) *v1.Instance {

	var lifecycleStatus v1.LifecycleStatus
	switch shadeformInstance.Status {
	case "creating":
		lifecycleStatus = v1.LifecycleStatusPending
	case "pending_provider":
		lifecycleStatus = v1.LifecycleStatusPending
	case "pending":
		lifecycleStatus = v1.LifecycleStatusPending
	case "active":
		lifecycleStatus = v1.LifecycleStatusRunning
	case "error":
		lifecycleStatus = v1.LifecycleStatusFailed
	default:
		lifecycleStatus = v1.LifecycleStatusPending
	}

	instanceType := c.getInstanceType(string(shadeformInstance.Cloud), shadeformInstance.ShadeInstanceType)

	instance := &v1.Instance{
		Name:         shadeformInstance.Name,
		CreatedAt:    shadeformInstance.CreatedAt,
		CloudID:      v1.CloudProviderInstanceID(shadeformInstance.Id),
		PublicIP:     shadeformInstance.Ip,
		Hostname:     hostname,
		ImageID:      shadeformInstance.Configuration.Os,
		InstanceType: instanceType,
		DiskSize:     units.Base2Bytes(shadeformInstance.Configuration.StorageInGb) * units.GiB,
		SSHUser:      shadeformInstance.SshUser,
		SSHPort:      int(shadeformInstance.SshPort),
		Status: v1.Status{
			LifecycleStatus: lifecycleStatus,
		},
		Spot:       false,
		Location:   shadeformInstance.Region,
		Stoppable:  false,
		Rebootable: true,
	}

	return instance
}
