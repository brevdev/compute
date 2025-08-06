package v1

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/units"
	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/compute/pkg/v1"
)

const lambdaLabsTimeNameFormat = "2006-01-02-15-04-05Z07-00"

// CreateInstance creates a new instance in Lambda Labs
// Supported via: POST /api/v1/instance-operations/launch
func (c *LambdaLabsClient) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
	keyPairName := attrs.RefID
	if attrs.KeyPairName != nil {
		keyPairName = *attrs.KeyPairName
	}

	if attrs.PublicKey != "" {
		err := c.addSSHKeyIdempotent(ctx, keyPairName, attrs.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	location := attrs.Location
	if location == "" {
		location = "us-west-1"
	}

	quantity := int32(1)
	request := openapi.LaunchInstanceRequest{
		RegionName:       location,
		InstanceTypeName: attrs.InstanceType,
		SshKeyNames:      []string{keyPairName},
		Quantity:         &quantity,
		FileSystemNames:  []string{},
	}

	name := fmt.Sprintf("%s--%s", c.GetReferenceID(), time.Now().UTC().Format(lambdaLabsTimeNameFormat))
	if attrs.Name != "" {
		name = fmt.Sprintf("%s--%s--%s", c.GetReferenceID(), attrs.Name, time.Now().UTC().Format(lambdaLabsTimeNameFormat))
	}
	request.Name = *openapi.NewNullableString(&name)

	resp, httpResp, err := c.client.DefaultAPI.LaunchInstance(c.makeAuthContext(ctx)).LaunchInstanceRequest(request).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, handleLLErrToCloudErr(err)
	}

	if len(resp.Data.InstanceIds) != 1 {
		return nil, fmt.Errorf("expected 1 instance ID, got %d", len(resp.Data.InstanceIds))
	}

	instanceID := v1.CloudProviderInstanceID(resp.Data.InstanceIds[0])
	return c.GetInstance(ctx, instanceID)
}

// GetInstance retrieves an instance by ID
// Supported via: GET /api/v1/instances/{id}
func (c *LambdaLabsClient) GetInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) (*v1.Instance, error) {
	resp, httpResp, err := c.client.DefaultAPI.GetInstance(c.makeAuthContext(ctx), string(instanceID)).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, handleLLErrToCloudErr(err)
	}

	return convertLambdaLabsInstanceToV1Instance(resp.Data), nil
}

// TerminateInstance terminates an instance
// Supported via: POST /api/v1/instance-operations/terminate
func (c *LambdaLabsClient) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	request := openapi.TerminateInstanceRequest{
		InstanceIds: []string{string(instanceID)},
	}

	_, httpResp, err := c.client.DefaultAPI.TerminateInstance(c.makeAuthContext(ctx)).TerminateInstanceRequest(request).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return handleLLErrToCloudErr(err)
	}

	return nil
}

// ListInstances lists all instances
// Supported via: GET /api/v1/instances
func (c *LambdaLabsClient) ListInstances(ctx context.Context, _ v1.ListInstancesArgs) ([]v1.Instance, error) {
	resp, httpResp, err := c.client.DefaultAPI.ListInstances(c.makeAuthContext(ctx)).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return nil, handleLLErrToCloudErr(err)
	}

	instances := make([]v1.Instance, 0, len(resp.Data))
	for _, llInstance := range resp.Data {
		instance := convertLambdaLabsInstanceToV1Instance(llInstance)
		instances = append(instances, *instance)
	}

	return instances, nil
}

func (c *LambdaLabsClient) addSSHKeyIdempotent(ctx context.Context, keyName, publicKey string) error {
	_, resp, err := c.client.DefaultAPI.AddSSHKey(c.makeAuthContext(ctx)).AddSSHKeyRequest(openapi.AddSSHKeyRequest{
		Name:      keyName,
		PublicKey: &publicKey,
	}).Execute()
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	if err != nil {
		if strings.Contains(err.Error(), "name must be unique") {
			return nil
		}
		return handleLLErrToCloudErr(err)
	}

	return nil
}

func parseGPUFromDescription(description string) v1.GPU {
	gpu := v1.GPU{
		Manufacturer: "NVIDIA",
	}
	
	countRegex := regexp.MustCompile(`(\d+)x`)
	if countMatch := countRegex.FindStringSubmatch(description); len(countMatch) > 1 {
		if count, err := strconv.ParseInt(countMatch[1], 10, 32); err == nil {
			gpu.Count = int32(count)
		}
	}
	
	memoryRegex := regexp.MustCompile(`\((\d+)\s*GB\)`)
	if memoryMatch := memoryRegex.FindStringSubmatch(description); len(memoryMatch) > 1 {
		if memoryGiB, err := strconv.Atoi(memoryMatch[1]); err == nil {
			gpu.Memory = units.Base2Bytes(memoryGiB) * units.GiB
		}
	}
	
	nameRegex := regexp.MustCompile(`\d+x\s+(.+?)\s+\(`)
	if nameMatch := nameRegex.FindStringSubmatch(description); len(nameMatch) > 1 {
		gpu.Name = strings.TrimSpace(nameMatch[1])
		gpu.Type = gpu.Name
	}
	
	if strings.Contains(description, "SXM4") {
		gpu.NetworkDetails = "SXM4"
	} else if strings.Contains(description, "PCIe") {
		gpu.NetworkDetails = "PCIe"
	}
	
	return gpu
}

func generateFirewallRuleFromPort(port int32) v1.FirewallRule {
	return v1.FirewallRule{
		FromPort: port,
		ToPort:   port,
		IPRanges: []string{"0.0.0.0/0"},
	}
}

// RebootInstance reboots an instance
// Supported via: POST /api/v1/instance-operations/restart
func (c *LambdaLabsClient) RebootInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
	request := openapi.RestartInstanceRequest{
		InstanceIds: []string{string(instanceID)},
	}

	_, httpResp, err := c.client.DefaultAPI.RestartInstance(c.makeAuthContext(ctx)).RestartInstanceRequest(request).Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	if err != nil {
		return handleLLErrToCloudErr(err)
	}

	return nil
}

// MergeInstanceForUpdate merges instance data for updates
func convertLambdaLabsInstanceToV1Instance(llInstance openapi.Instance) *v1.Instance {
	var publicIP, privateIP, hostname, name string

	if llInstance.Ip.IsSet() {
		publicIP = *llInstance.Ip.Get()
	}
	if llInstance.PrivateIp.IsSet() {
		privateIP = *llInstance.PrivateIp.Get()
	}
	if llInstance.Hostname.IsSet() {
		hostname = *llInstance.Hostname.Get()
	}
	if llInstance.Name.IsSet() {
		name = *llInstance.Name.Get()
	}

	var cloudCredRefID string
	var createdAt time.Time
	if name != "" {
		parts := strings.Split(name, "--")
		if len(parts) > 0 {
			cloudCredRefID = parts[0]
		}
		if len(parts) > 1 {
			createdAt, _ = time.Parse("2006-01-02-15-04-05Z07-00", parts[1])
		}
	}

	refID := ""
	if len(llInstance.SshKeyNames) > 0 {
		refID = llInstance.SshKeyNames[0]
	}

	firewallRules := v1.FirewallRules{
		IngressRules: []v1.FirewallRule{
			generateFirewallRuleFromPort(22),
			generateFirewallRuleFromPort(2222),
		},
		EgressRules: []v1.FirewallRule{
			generateFirewallRuleFromPort(22),
			generateFirewallRuleFromPort(2222),
		},
	}

	return &v1.Instance{
		Name:           name,
		RefID:          refID,
		CloudCredRefID: cloudCredRefID,
		CreatedAt:      createdAt,
		CloudID:        v1.CloudProviderInstanceID(llInstance.Id),
		PublicIP:       publicIP,
		PrivateIP:      privateIP,
		PublicDNS:      publicIP,
		Hostname:       hostname,
		InstanceType:   llInstance.InstanceType.Name,
		Status: v1.Status{
			LifecycleStatus: convertLambdaLabsStatusToV1Status(llInstance.Status),
		},
		Location:      llInstance.Region.Name,
		SSHUser:       "ubuntu",
		SSHPort:       22,
		Stoppable:     false,
		Rebootable:    true,
		VolumeType:    "ssd",
		DiskSize:      units.Base2Bytes(llInstance.InstanceType.Specs.StorageGib) * units.GiB,
		FirewallRules: firewallRules,
	}
}

func convertLambdaLabsStatusToV1Status(status string) v1.LifecycleStatus {
	switch status {
	case "booting":
		return v1.LifecycleStatusPending
	case "active":
		return v1.LifecycleStatusRunning
	case "terminating":
		return v1.LifecycleStatusTerminating
	case "terminated":
		return v1.LifecycleStatusTerminated
	case "error":
		return v1.LifecycleStatusFailed
	case "unhealthy":
		return v1.LifecycleStatusRunning
	default:
		return v1.LifecycleStatusPending
	}
}

func (c *LambdaLabsClient) MergeInstanceForUpdate(_ v1.Instance, newInst v1.Instance) v1.Instance {
	return newInst
}

func (c *LambdaLabsClient) MergeInstanceTypeForUpdate(_ v1.InstanceType, newIt v1.InstanceType) v1.InstanceType {
	return newIt
}
