package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/compute/pkg/v1"
)

const (
	testInstanceID      = "test-instance-id"
	nonexistentInstance = "nonexistent-instance"
)

func setupMockClient() (*LambdaLabsClient, func()) {
	httpmock.Activate()
	client := NewLambdaLabsClient("test-ref-id", "test-api-key")
	return client, httpmock.DeactivateAndReset
}

func TestLambdaLabsCredential_GetReferenceID(t *testing.T) {
	cred := &LambdaLabsCredential{
		RefID:  "test-ref-id",
		APIKey: "test-api-key",
	}

	assert.Equal(t, "test-ref-id", cred.GetReferenceID())
}

func TestLambdaLabsCredential_GetAPIType(t *testing.T) {
	cred := &LambdaLabsCredential{}
	assert.Equal(t, v1.APITypeGlobal, cred.GetAPIType())
}

func TestLambdaLabsCredential_GetCloudProviderID(t *testing.T) {
	cred := &LambdaLabsCredential{}
	assert.Equal(t, v1.CloudProviderID("lambdalabs"), cred.GetCloudProviderID())
}

func TestLambdaLabsCredential_GetTenantID(t *testing.T) {
	cred := &LambdaLabsCredential{APIKey: "test-key"}
	tenantID, err := cred.GetTenantID()
	assert.NoError(t, err)
	assert.Contains(t, tenantID, "lambdalabs-")
}

func TestLambdaLabsCredential_MakeClient(t *testing.T) {
	cred := &LambdaLabsCredential{
		RefID:  "test-ref-id",
		APIKey: "test-api-key",
	}

	client, err := cred.MakeClient(context.Background(), "test-tenant")
	require.NoError(t, err)
	lambdaClient, ok := client.(*LambdaLabsClient)
	require.True(t, ok)
	assert.Equal(t, "test-ref-id", lambdaClient.refID)
	assert.Equal(t, "test-api-key", lambdaClient.apiKey)
}

func TestLambdaLabsClient_GetAPIType(t *testing.T) {
	client := &LambdaLabsClient{}
	assert.Equal(t, v1.APITypeGlobal, client.GetAPIType())
}

func TestLambdaLabsClient_GetCloudProviderID(t *testing.T) {
	client := &LambdaLabsClient{}
	assert.Equal(t, v1.CloudProviderID("lambdalabs"), client.GetCloudProviderID())
}

func TestLambdaLabsClient_MakeClient(t *testing.T) {
	client := &LambdaLabsClient{
		refID:  "test-ref-id",
		apiKey: "test-api-key",
	}

	newClient, err := client.MakeClient(context.Background(), "test-tenant")
	require.NoError(t, err)
	lambdaClient, ok := newClient.(*LambdaLabsClient)
	require.True(t, ok)
	assert.Equal(t, client, lambdaClient)
}

func TestLambdaLabsClient_GetReferenceID(t *testing.T) {
	client := &LambdaLabsClient{refID: "test-ref-id"}
	assert.Equal(t, "test-ref-id", client.GetReferenceID())
}

func TestLambdaLabsClient_makeAuthContext(t *testing.T) {
	client := &LambdaLabsClient{apiKey: "test-api-key"}
	ctx := client.makeAuthContext(context.Background())

	auth := ctx.Value(openapi.ContextBasicAuth)
	require.NotNil(t, auth)

	basicAuth, ok := auth.(openapi.BasicAuth)
	require.True(t, ok)
	assert.Equal(t, "test-api-key", basicAuth.UserName)
	assert.Equal(t, "", basicAuth.Password)
}

func TestLambdaLabsClient_CreateInstance_Success(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := testInstanceID
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ test@example.com"

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/ssh-keys",
		httpmock.NewJsonResponderOrPanic(200, openapi.AddSSHKey200Response{
			Data: openapi.SshKey{
				Id:        "ssh-key-id",
				Name:      "test-instance-id",
				PublicKey: publicKey,
			},
		}))

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/launch",
		httpmock.NewJsonResponderOrPanic(200, openapi.LaunchInstance200Response{
			Data: openapi.LaunchInstance200ResponseData{
				InstanceIds: []string{instanceID},
			},
		}))

	mockInstance := createMockInstance(instanceID)
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cloud.lambda.ai/api/v1/instances/%s", instanceID),
		httpmock.NewJsonResponderOrPanic(200, openapi.GetInstance200Response{
			Data: mockInstance,
		}))

	args := v1.CreateInstanceAttrs{
		InstanceType: "gpu_1x_a10",
		Location:     "us-west-1",
		PublicKey:    publicKey,
		Name:         "test-instance",
	}

	instance, err := client.CreateInstance(context.Background(), args)
	require.NoError(t, err)
	assert.Equal(t, instanceID, string(instance.CloudID))
	assert.Contains(t, instance.Name, "test-instance")
	assert.Equal(t, v1.LifecycleStatusRunning, instance.Status.LifecycleStatus)
}

func TestLambdaLabsClient_CreateInstance_WithoutPublicKey(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := testInstanceID

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/launch",
		httpmock.NewJsonResponderOrPanic(200, openapi.LaunchInstance200Response{
			Data: openapi.LaunchInstance200ResponseData{
				InstanceIds: []string{instanceID},
			},
		}))

	mockInstance := createMockInstance(instanceID)
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cloud.lambda.ai/api/v1/instances/%s", instanceID),
		httpmock.NewJsonResponderOrPanic(200, openapi.GetInstance200Response{
			Data: mockInstance,
		}))

	args := v1.CreateInstanceAttrs{
		InstanceType: "gpu_1x_a10",
		Location:     "us-west-1",
		Name:         "test-instance",
	}

	instance, err := client.CreateInstance(context.Background(), args)
	require.NoError(t, err)
	assert.Equal(t, instanceID, string(instance.CloudID))
}

func TestLambdaLabsClient_CreateInstance_SSHKeyError(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ test@example.com"

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/ssh-keys",
		httpmock.NewStringResponder(400, `{"error": {"code": "INVALID_REQUEST", "message": "SSH key already exists"}}`))

	args := v1.CreateInstanceAttrs{
		InstanceType: "gpu_1x_a10",
		Location:     "us-west-1",
		PublicKey:    publicKey,
		Name:         "test-instance",
	}

	_, err := client.CreateInstance(context.Background(), args)
	assert.Error(t, err)
}

func TestLambdaLabsClient_CreateInstance_LaunchError(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/launch",
		httpmock.NewStringResponder(400, `{"error": {"code": "INVALID_REQUEST", "message": "Instance type not available"}}`))

	args := v1.CreateInstanceAttrs{
		InstanceType: "gpu_1x_a10",
		Location:     "us-west-1",
		Name:         "test-instance",
	}

	_, err := client.CreateInstance(context.Background(), args)
	assert.Error(t, err)
}

func TestLambdaLabsClient_GetInstance_Success(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := testInstanceID
	mockInstance := createMockInstance(instanceID)

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cloud.lambda.ai/api/v1/instances/%s", instanceID),
		httpmock.NewJsonResponderOrPanic(200, openapi.GetInstance200Response{
			Data: mockInstance,
		}))

	instance, err := client.GetInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	require.NoError(t, err)
	assert.Equal(t, instanceID, string(instance.CloudID))
	assert.Equal(t, "test-instance", instance.Name)
	assert.Equal(t, v1.LifecycleStatusRunning, instance.Status.LifecycleStatus)
}

func TestLambdaLabsClient_GetInstance_NotFound(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := nonexistentInstance

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://cloud.lambda.ai/api/v1/instances/%s", instanceID),
		httpmock.NewStringResponder(404, `{"error": {"code": "NOT_FOUND", "message": "Instance not found"}}`))

	_, err := client.GetInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	assert.Error(t, err)
}

func TestLambdaLabsClient_ListInstances_Success(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	mockInstances := []openapi.Instance{
		createMockInstance("instance-1"),
		createMockInstance("instance-2"),
	}

	httpmock.RegisterResponder("GET", "https://cloud.lambda.ai/api/v1/instances",
		httpmock.NewJsonResponderOrPanic(200, openapi.ListInstances200Response{
			Data: mockInstances,
		}))

	instances, err := client.ListInstances(context.Background(), v1.ListInstancesArgs{})
	require.NoError(t, err)
	assert.Len(t, instances, 2)
	assert.Equal(t, "instance-1", string(instances[0].CloudID))
	assert.Equal(t, "instance-2", string(instances[1].CloudID))
}

func TestLambdaLabsClient_ListInstances_Empty(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://cloud.lambda.ai/api/v1/instances",
		httpmock.NewJsonResponderOrPanic(200, openapi.ListInstances200Response{
			Data: []openapi.Instance{},
		}))

	instances, err := client.ListInstances(context.Background(), v1.ListInstancesArgs{})
	require.NoError(t, err)
	assert.Len(t, instances, 0)
}

func TestLambdaLabsClient_ListInstances_Error(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://cloud.lambda.ai/api/v1/instances",
		httpmock.NewStringResponder(500, `{"error": {"code": "INTERNAL_ERROR", "message": "Internal server error"}}`))

	_, err := client.ListInstances(context.Background(), v1.ListInstancesArgs{})
	assert.Error(t, err)
}

func TestLambdaLabsClient_TerminateInstance_Success(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := testInstanceID

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/terminate",
		httpmock.NewJsonResponderOrPanic(200, openapi.TerminateInstance200Response{
			Data: openapi.TerminateInstance200ResponseData{
				TerminatedInstances: []openapi.Instance{
					createMockInstance(instanceID),
				},
			},
		}))

	err := client.TerminateInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	assert.NoError(t, err)
}

func TestLambdaLabsClient_TerminateInstance_Error(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := nonexistentInstance

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/terminate",
		httpmock.NewStringResponder(404, `{"error": {"code": "NOT_FOUND", "message": "Instance not found"}}`))

	err := client.TerminateInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	assert.Error(t, err)
}

func TestLambdaLabsClient_RebootInstance_Success(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := testInstanceID

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/restart",
		httpmock.NewJsonResponderOrPanic(200, openapi.RestartInstance200Response{
			Data: openapi.RestartInstance200ResponseData{
				RestartedInstances: []openapi.Instance{
					createMockInstance(instanceID),
				},
			},
		}))

	err := client.RebootInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	assert.NoError(t, err)
}

func TestLambdaLabsClient_RebootInstance_Error(t *testing.T) {
	client, cleanup := setupMockClient()
	defer cleanup()

	instanceID := nonexistentInstance

	httpmock.RegisterResponder("POST", "https://cloud.lambda.ai/api/v1/instance-operations/restart",
		httpmock.NewStringResponder(404, `{"error": {"code": "NOT_FOUND", "message": "Instance not found"}}`))

	err := client.RebootInstance(context.Background(), v1.CloudProviderInstanceID(instanceID))
	assert.Error(t, err)
}

func TestLambdaLabsClient_GetCapabilities(t *testing.T) {
	client := &LambdaLabsClient{}
	capabilities, err := client.GetCapabilities(context.Background())
	require.NoError(t, err)

	assert.Contains(t, capabilities, v1.CapabilityCreateInstance)
	assert.Contains(t, capabilities, v1.CapabilityTerminateInstance)
	assert.Contains(t, capabilities, v1.CapabilityRebootInstance)
	assert.NotContains(t, capabilities, v1.CapabilityStopStartInstance)
}

func TestLambdaLabsCredential_GetCapabilities(t *testing.T) {
	cred := &LambdaLabsCredential{}
	capabilities, err := cred.GetCapabilities(context.Background())
	require.NoError(t, err)

	assert.Contains(t, capabilities, v1.CapabilityCreateInstance)
	assert.Contains(t, capabilities, v1.CapabilityTerminateInstance)
	assert.Contains(t, capabilities, v1.CapabilityRebootInstance)
	assert.NotContains(t, capabilities, v1.CapabilityStopStartInstance)
}

func TestConvertLambdaLabsInstanceToV1Instance(t *testing.T) {
	lambdaInstance := createMockInstance("test-instance-id")

	v1Instance := convertLambdaLabsInstanceToV1Instance(lambdaInstance)

	assert.Equal(t, "test-instance-id", string(v1Instance.CloudID))
	assert.Equal(t, "test-instance", v1Instance.Name)
	assert.Equal(t, v1.LifecycleStatusRunning, v1Instance.Status.LifecycleStatus)
	assert.Equal(t, "192.168.1.100", v1Instance.PublicIP)
	assert.Equal(t, "10.0.1.100", v1Instance.PrivateIP)
	assert.Equal(t, "us-west-1", v1Instance.Location)
	assert.Equal(t, "gpu_1x_a10", v1Instance.InstanceType)
}

func TestConvertLambdaLabsStatusToV1Status(t *testing.T) {
	tests := []struct {
		lambdaStatus string
		expected     v1.LifecycleStatus
	}{
		{"active", v1.LifecycleStatusRunning},
		{"booting", v1.LifecycleStatusPending},
		{"unhealthy", v1.LifecycleStatusRunning},
		{"terminating", v1.LifecycleStatusTerminating},
		{"terminated", v1.LifecycleStatusTerminated},
		{"error", v1.LifecycleStatusFailed},
	}

	for _, test := range tests {
		t.Run(test.lambdaStatus, func(t *testing.T) {
			result := convertLambdaLabsStatusToV1Status(test.lambdaStatus)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMergeInstanceForUpdate(t *testing.T) {
	client := &LambdaLabsClient{}
	original := v1.Instance{
		CloudID: "test-id",
		Name:    "original-name",
		Status:  v1.Status{LifecycleStatus: v1.LifecycleStatusRunning},
	}

	update := v1.Instance{
		Name:   "updated-name",
		Status: v1.Status{LifecycleStatus: v1.LifecycleStatusTerminated},
	}

	merged := client.MergeInstanceForUpdate(original, update)

	assert.Equal(t, "updated-name", merged.Name)
	assert.Equal(t, v1.LifecycleStatusTerminated, merged.Status.LifecycleStatus)
}

func TestMergeInstanceTypeForUpdate(t *testing.T) {
	client := &LambdaLabsClient{}
	original := v1.InstanceType{
		ID:   "test-id",
		Type: "original-type",
	}

	update := v1.InstanceType{
		Type: "updated-type",
	}

	merged := client.MergeInstanceTypeForUpdate(original, update)

	assert.Equal(t, "updated-type", merged.Type)
}

func createMockInstance(instanceID string) openapi.Instance {
	name := "test-instance"
	ip := "192.168.1.100"
	privateIP := "10.0.1.100"
	hostname := "test-instance.lambda.ai"

	return openapi.Instance{
		Id:              instanceID,
		Name:            *openapi.NewNullableString(&name),
		Ip:              *openapi.NewNullableString(&ip),
		PrivateIp:       *openapi.NewNullableString(&privateIP),
		Status:          "active",
		SshKeyNames:     []string{"test-key"},
		FileSystemNames: []string{},
		Region: &openapi.Region{
			Name:        "us-west-1",
			Description: "US West 1",
		},
		InstanceType: &openapi.InstanceType{
			Name:              "gpu_1x_a10",
			Description:       "1x NVIDIA A10 GPU",
			GpuDescription:    "NVIDIA A10",
			PriceCentsPerHour: 100,
			Specs: openapi.InstanceTypeSpecs{
				MemoryGib:  32,
				StorageGib: 512,
			},
		},
		Hostname: *openapi.NewNullableString(&hostname),
	}
}
