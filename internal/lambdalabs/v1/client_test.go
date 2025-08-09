package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	openapi "github.com/brevdev/cloud/internal/lambdalabs/gen/lambdalabs"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

func TestLambdaLabsClient_GetAPIType(t *testing.T) {
	client := &LambdaLabsClient{}
	assert.Equal(t, v1.APITypeGlobal, client.GetAPIType())
}

func TestLambdaLabsClient_GetCloudProviderID(t *testing.T) {
	client := &LambdaLabsClient{}
	assert.Equal(t, v1.CloudProviderID("lambda-labs"), client.GetCloudProviderID())
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
