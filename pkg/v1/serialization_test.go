package v1

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockCredential struct {
	providerID CloudProviderID
	refID      string
	tenantID   string
}

func (m *MockCredential) MakeClient(_ context.Context, _ string) (CloudClient, error) {
	return nil, ErrNotImplemented
}

func (m *MockCredential) GetTenantID() (string, error) {
	return m.tenantID, nil
}

func (m *MockCredential) GetReferenceID() string {
	return m.refID
}

func (m *MockCredential) GetAPIType() APIType {
	return APITypeGlobal
}

func (m *MockCredential) GetCapabilities(_ context.Context) (Capabilities, error) {
	return nil, ErrNotImplemented
}

func (m *MockCredential) GetCloudProviderID() CloudProviderID {
	return m.providerID
}

func TestSerializedCredential_Structure(t *testing.T) {
	serialized := SerializedCredential{
		ProviderID: "test-provider",
		Data:       json.RawMessage(`{"test": "data"}`),
	}

	bytes, err := json.Marshal(serialized)
	require.NoError(t, err)

	var unmarshaled SerializedCredential
	err = json.Unmarshal(bytes, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, "test-provider", unmarshaled.ProviderID)
	assert.Equal(t, json.RawMessage(`{"test":"data"}`), unmarshaled.Data)
}

func TestCredentialDataStructures_JSONTags(t *testing.T) {
	t.Run("LambdaLabsCredentialData", func(t *testing.T) {
		data := LambdaLabsCredentialData{
			RefID:  "test-ref",
			APIKey: "test-key",
		}

		bytes, err := json.Marshal(data)
		require.NoError(t, err)

		expected := `{"ref_id":"test-ref","api_key":"test-key"}`
		assert.JSONEq(t, expected, string(bytes))

		var unmarshaled LambdaLabsCredentialData
		err = json.Unmarshal(bytes, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, data, unmarshaled)
	})

	t.Run("FluidStackCredentialData", func(t *testing.T) {
		data := FluidStackCredentialData{
			RefID:  "test-ref",
			APIKey: "test-key",
		}

		bytes, err := json.Marshal(data)
		require.NoError(t, err)

		expected := `{"ref_id":"test-ref","api_key":"test-key"}`
		assert.JSONEq(t, expected, string(bytes))

		var unmarshaled FluidStackCredentialData
		err = json.Unmarshal(bytes, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, data, unmarshaled)
	})

	t.Run("NebiusCredentialData", func(t *testing.T) {
		data := NebiusCredentialData{
			RefID:             "test-ref",
			ServiceAccountKey: "test-key",
			ProjectID:         "test-project",
		}

		bytes, err := json.Marshal(data)
		require.NoError(t, err)

		expected := `{"ref_id":"test-ref","service_account_key":"test-key","project_id":"test-project"}`
		assert.JSONEq(t, expected, string(bytes))

		var unmarshaled NebiusCredentialData
		err = json.Unmarshal(bytes, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, data, unmarshaled)
	})
}

type MockSerializableCredential struct {
	MockCredential
	data interface{}
}

func (m *MockSerializableCredential) SerializeData() (interface{}, error) {
	return m.data, nil
}

func TestSerializeCredentialData(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		credData := LambdaLabsCredentialData{
			RefID:  "test-ref",
			APIKey: "test-key",
		}

		bytes, err := SerializeCredentialData("lambda-labs", credData)
		require.NoError(t, err)

		var serialized SerializedCredential
		err = json.Unmarshal(bytes, &serialized)
		require.NoError(t, err)

		assert.Equal(t, "lambda-labs", serialized.ProviderID)

		var unmarshaled LambdaLabsCredentialData
		err = json.Unmarshal(serialized.Data, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, credData, unmarshaled)
	})

	t.Run("empty provider ID", func(t *testing.T) {
		_, err := SerializeCredentialData("", map[string]string{"test": "data"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provider_id cannot be empty")
	})
}

func TestSerializeCredentialDataToString(t *testing.T) {
	credData := FluidStackCredentialData{
		RefID:  "test-ref",
		APIKey: "test-key",
	}

	str, err := SerializeCredentialDataToString("fluidstack", credData)
	require.NoError(t, err)
	assert.NotEmpty(t, str)

	var serialized SerializedCredential
	err = json.Unmarshal([]byte(str), &serialized)
	require.NoError(t, err)
	assert.Equal(t, "fluidstack", serialized.ProviderID)
}

func TestSerializeCredential_ErrorCases(t *testing.T) {
	t.Run("nil credential", func(t *testing.T) {
		_, err := SerializeCredential(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "credential cannot be nil")
	})

	t.Run("non-serializable credential", func(t *testing.T) {
		cred := &MockCredential{
			providerID: "lambda-labs",
			refID:      "test-ref",
		}

		_, err := SerializeCredential(cred)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not implement SerializableCredential")
	})

	t.Run("empty provider ID", func(t *testing.T) {
		cred := &MockSerializableCredential{
			MockCredential: MockCredential{
				providerID: "",
				refID:      "test-ref",
			},
			data: map[string]string{"test": "data"},
		}

		_, err := SerializeCredential(cred)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "credential must have a valid provider ID")
	})

	t.Run("valid serializable credential", func(t *testing.T) {
		cred := &MockSerializableCredential{
			MockCredential: MockCredential{
				providerID: "lambda-labs",
				refID:      "test-ref",
			},
			data: LambdaLabsCredentialData{
				RefID:  "test-ref",
				APIKey: "test-key",
			},
		}

		bytes, err := SerializeCredential(cred)
		require.NoError(t, err)
		assert.NotEmpty(t, bytes)

		var serialized SerializedCredential
		err = json.Unmarshal(bytes, &serialized)
		require.NoError(t, err)
		assert.Equal(t, "lambda-labs", serialized.ProviderID)
	})
}

func TestDeserializeCredential_ErrorCases(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		_, err := DeserializeCredential([]byte{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data cannot be empty")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := DeserializeCredential([]byte("invalid json"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal serialized credential")
	})

	t.Run("empty provider ID", func(t *testing.T) {
		serialized := SerializedCredential{
			ProviderID: "",
			Data:       json.RawMessage(`{"test": "data"}`),
		}
		bytes, _ := json.Marshal(serialized)

		_, err := DeserializeCredential(bytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provider_id cannot be empty")
	})

	t.Run("unsupported provider", func(t *testing.T) {
		serialized := SerializedCredential{
			ProviderID: "unsupported-provider",
			Data:       json.RawMessage(`{"test": "data"}`),
		}
		bytes, _ := json.Marshal(serialized)

		_, err := DeserializeCredential(bytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported provider: unsupported-provider")
	})

	t.Run("supported provider but not implemented", func(t *testing.T) {
		serialized := SerializedCredential{
			ProviderID: "lambda-labs",
			Data:       json.RawMessage(`{"ref_id": "test", "api_key": "key"}`),
		}
		bytes, _ := json.Marshal(serialized)

		_, err := DeserializeCredential(bytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires import")
	})
}

func TestDeserializeCredentialFromString(t *testing.T) {
	_, err := DeserializeCredentialFromString("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be empty")
}

func TestDeserializeCredentialByProvider(t *testing.T) {
	t.Run("lambda-labs with valid data", func(t *testing.T) {
		data := json.RawMessage(`{"ref_id": "test-ref", "api_key": "test-key"}`)
		_, err := DeserializeCredentialByProvider("lambda-labs", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires import")
	})

	t.Run("lambda-labs with missing ref_id", func(t *testing.T) {
		data := json.RawMessage(`{"api_key": "test-key"}`)
		_, err := DeserializeCredentialByProvider("lambda-labs", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have a reference ID")
	})

	t.Run("lambda-labs with missing api_key", func(t *testing.T) {
		data := json.RawMessage(`{"ref_id": "test-ref"}`)
		_, err := DeserializeCredentialByProvider("lambda-labs", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have an API key")
	})

	t.Run("fluidstack with valid data", func(t *testing.T) {
		data := json.RawMessage(`{"ref_id": "test-ref", "api_key": "test-key"}`)
		_, err := DeserializeCredentialByProvider("fluidstack", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires import")
	})

	t.Run("nebius with valid data", func(t *testing.T) {
		data := json.RawMessage(`{"ref_id": "test-ref", "service_account_key": "test-key", "project_id": "test-project"}`)
		_, err := DeserializeCredentialByProvider("nebius", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires import")
	})

	t.Run("nebius with missing project_id", func(t *testing.T) {
		data := json.RawMessage(`{"ref_id": "test-ref", "service_account_key": "test-key"}`)
		_, err := DeserializeCredentialByProvider("nebius", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have a project ID")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := json.RawMessage(`invalid json`)
		_, err := DeserializeCredentialByProvider("lambda-labs", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal")
	})
}

func TestRoundTripSerialization(t *testing.T) {
	t.Run("credential data structures", func(t *testing.T) {
		testCases := []struct {
			name       string
			providerID string
			data       interface{}
		}{
			{
				name:       "lambda-labs",
				providerID: "lambda-labs",
				data: LambdaLabsCredentialData{
					RefID:  "test-ref",
					APIKey: "test-key",
				},
			},
			{
				name:       "fluidstack",
				providerID: "fluidstack",
				data: FluidStackCredentialData{
					RefID:  "test-ref",
					APIKey: "test-key",
				},
			},
			{
				name:       "nebius",
				providerID: "nebius",
				data: NebiusCredentialData{
					RefID:             "test-ref",
					ServiceAccountKey: "test-key",
					ProjectID:         "test-project",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				serialized, err := SerializeCredentialData(tc.providerID, tc.data)
				require.NoError(t, err)

				var wrapper SerializedCredential
				err = json.Unmarshal(serialized, &wrapper)
				require.NoError(t, err)

				assert.Equal(t, tc.providerID, wrapper.ProviderID)

				switch tc.providerID {
				case "lambda-labs":
					var data LambdaLabsCredentialData
					err = json.Unmarshal(wrapper.Data, &data)
					require.NoError(t, err)
					assert.Equal(t, tc.data, data)
				case "fluidstack":
					var data FluidStackCredentialData
					err = json.Unmarshal(wrapper.Data, &data)
					require.NoError(t, err)
					assert.Equal(t, tc.data, data)
				case "nebius":
					var data NebiusCredentialData
					err = json.Unmarshal(wrapper.Data, &data)
					require.NoError(t, err)
					assert.Equal(t, tc.data, data)
				}
			})
		}
	})
}
