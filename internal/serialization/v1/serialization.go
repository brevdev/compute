package v1

import (
	"encoding/json"
	"fmt"

	fluidstackv1 "github.com/brevdev/cloud/internal/fluidstack/v1"
	lambdalabsv1 "github.com/brevdev/cloud/internal/lambdalabs/v1"
	nebiusv1 "github.com/brevdev/cloud/internal/nebius/v1"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

type SerializedCredential struct {
	ProviderID string          `json:"provider_id"`
	Data       json.RawMessage `json:"data"`
}

type SerializableCredential interface {
	v1.CloudCredential
	SerializeData() (interface{}, error)
}

func SerializeCredentialData(providerID string, credData interface{}) ([]byte, error) {
	if providerID == "" {
		return nil, fmt.Errorf("provider_id cannot be empty")
	}

	dataBytes, err := json.Marshal(credData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credential data: %w", err)
	}

	serialized := SerializedCredential{
		ProviderID: providerID,
		Data:       dataBytes,
	}

	return json.Marshal(serialized)
}

func SerializeCredentialDataToString(providerID string, credData interface{}) (string, error) {
	bytes, err := SerializeCredentialData(providerID, credData)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func SerializeCredentialDataToJSON(providerID string, credData interface{}) ([]byte, error) {
	return SerializeCredentialData(providerID, credData)
}

func SerializeCredential(cred v1.CloudCredential) ([]byte, error) {
	if cred == nil {
		return nil, fmt.Errorf("credential cannot be nil")
	}

	serializableCred, ok := cred.(SerializableCredential)
	if !ok {
		return nil, fmt.Errorf("credential does not implement SerializableCredential interface")
	}

	providerID := string(cred.GetCloudProviderID())
	if providerID == "" {
		return nil, fmt.Errorf("credential must have a valid provider ID")
	}

	credData, err := serializableCred.SerializeData()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize credential data: %w", err)
	}

	return SerializeCredentialData(providerID, credData)
}

func SerializeCredentialToString(cred v1.CloudCredential) (string, error) {
	bytes, err := SerializeCredential(cred)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func SerializeCredentialToJSON(cred v1.CloudCredential) ([]byte, error) {
	return SerializeCredential(cred)
}

func DeserializeCredential(data []byte) (v1.CloudCredential, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	var serialized SerializedCredential
	if err := json.Unmarshal(data, &serialized); err != nil {
		return nil, fmt.Errorf("failed to unmarshal serialized credential: %w", err)
	}

	if serialized.ProviderID == "" {
		return nil, fmt.Errorf("provider_id cannot be empty")
	}

	return DeserializeCredentialByProvider(serialized.ProviderID, serialized.Data)
}

func DeserializeCredentialFromString(data string) (v1.CloudCredential, error) {
	return DeserializeCredential([]byte(data))
}

func DeserializeCredentialFromJSON(data []byte) (v1.CloudCredential, error) {
	return DeserializeCredential(data)
}

func DeserializeCredentialByProvider(providerID string, data json.RawMessage) (v1.CloudCredential, error) {
	switch providerID {
	case lambdalabsv1.CloudProviderID:
		return DeserializeLambdaLabsCredential(data)
	case fluidstackv1.CloudProviderID:
		return DeserializeFluidStackCredential(data)
	case "nebius":
		return DeserializeNebiusCredential(data)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerID)
	}
}

func DeserializeLambdaLabsCredential(data json.RawMessage) (v1.CloudCredential, error) {
	var credData lambdalabsv1.LambdaLabsCredential
	if err := json.Unmarshal(data, &credData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lambda labs credential data: %w", err)
	}

	if credData.RefID == "" {
		return nil, fmt.Errorf("lambda labs credential must have a reference ID")
	}

	if credData.APIKey == "" {
		return nil, fmt.Errorf("lambda labs credential must have an API key")
	}

	return &credData, nil
}

func DeserializeFluidStackCredential(data json.RawMessage) (v1.CloudCredential, error) {
	var credData fluidstackv1.FluidStackCredential
	if err := json.Unmarshal(data, &credData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fluidstack credential data: %w", err)
	}

	if credData.RefID == "" {
		return nil, fmt.Errorf("fluidstack credential must have a reference ID")
	}

	if credData.APIKey == "" {
		return nil, fmt.Errorf("fluidstack credential must have an API key")
	}

	return &credData, nil
}

func DeserializeNebiusCredential(data json.RawMessage) (v1.CloudCredential, error) {
	var credData nebiusv1.NebiusCredential
	if err := json.Unmarshal(data, &credData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nebius credential data: %w", err)
	}

	if credData.RefID == "" {
		return nil, fmt.Errorf("nebius credential must have a reference ID")
	}

	if credData.ServiceAccountKey == "" {
		return nil, fmt.Errorf("nebius credential must have a service account key")
	}

	if credData.ProjectID == "" {
		return nil, fmt.Errorf("nebius credential must have a project ID")
	}

	return &credData, nil
}
