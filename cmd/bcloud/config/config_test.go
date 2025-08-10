package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	lambdalabs "github.com/brevdev/cloud/internal/lambdalabs/v1"
	"gopkg.in/yaml.v3"
)

func TestRegisteredProvidersHaveRefIDField(t *testing.T) {
	for providerID, factory := range providerRegistry {
		t.Run(providerID, func(t *testing.T) {
			cred := factory()
			credType := reflect.TypeOf(cred)

			if credType.Kind() == reflect.Ptr {
				credType = credType.Elem()
			}

			hasRefID := false
			for i := 0; i < credType.NumField(); i++ {
				field := credType.Field(i)
				if field.Name == "RefID" {
					hasRefID = true
					break
				}
				if field.Anonymous && field.Type.Kind() == reflect.Ptr {
					embeddedType := field.Type.Elem()
					for j := 0; j < embeddedType.NumField(); j++ {
						embeddedField := embeddedType.Field(j)
						if embeddedField.Name == "RefID" {
							hasRefID = true
							break
						}
					}
				}
			}

			if !hasRefID {
				t.Errorf("Provider %s does not have a RefID field in its credential struct", providerID)
			}
		})
	}
}

func TestLoadConfig_Success(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".bcloud", "credentials.yaml")

	err := os.MkdirAll(filepath.Dir(configPath), 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `
credentials:
  test-lambda:
    provider: "lambdalabs"
    APIKey: "test-key"
    default_location: "us-west-1"
  test-lambda-with-ref:
    provider: "lambdalabs"
    APIKey: "test-key-2"
    RefID: "custom-ref-id"
    default_location: "us-east-1"
settings:
  output_format: "json"
  default_timeout: "10m"
`

	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	_ = os.Setenv("HOME", tempDir)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(config.Credentials) != 2 {
		t.Errorf("Expected 2 credentials, got %d", len(config.Credentials))
	}

	testLambda, exists := config.Credentials["test-lambda"]
	if !exists {
		t.Fatal("test-lambda credential not found")
	}
	if testLambda.Provider != "lambdalabs" {
		t.Errorf("Expected provider 'lambdalabs', got '%s'", testLambda.Provider)
	}
	if testLambda.Value.GetReferenceID() != "test-lambda" {
		t.Errorf("Expected ref_id 'test-lambda', got '%s'", testLambda.Value.GetReferenceID())
	}

	testLambdaWithRef, exists := config.Credentials["test-lambda-with-ref"]
	if !exists {
		t.Fatal("test-lambda-with-ref credential not found")
	}
	if testLambdaWithRef.Value.GetReferenceID() != "custom-ref-id" {
		t.Errorf("Expected ref_id 'custom-ref-id', got '%s'", testLambdaWithRef.Value.GetReferenceID())
	}

	if config.Settings.OutputFormat != "json" {
		t.Errorf("Expected output_format 'json', got '%s'", config.Settings.OutputFormat)
	}
	if config.Settings.DefaultTimeout != "10m" {
		t.Errorf("Expected default_timeout '10m', got '%s'", config.Settings.DefaultTimeout)
	}
}

func TestLoadConfig_WorkingDirectoryFallback(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(tempDir)

	originalHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	_ = os.Setenv("HOME", "/nonexistent")

	configContent := `
credentials:
  test-cred:
    provider: "lambdalabs"
    APIKey: "test-key"
    default_location: "us-west-1"
`

	err := os.WriteFile("bcloud.yaml", []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(config.Credentials) != 1 {
		t.Errorf("Expected 1 credential, got %d", len(config.Credentials))
	}
}

func TestLoadConfig_NoConfigFound(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	_ = os.Setenv("HOME", "/nonexistent")

	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	tempDir := t.TempDir()
	_ = os.Chdir(tempDir)

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error when no config found")
	}
	if !strings.Contains(err.Error(), "no configuration found") {
		t.Errorf("Expected 'no configuration found' error, got: %v", err)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".bcloud", "credentials.yaml")

	err := os.MkdirAll(filepath.Dir(configPath), 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	invalidYAML := `
credentials:
  test-cred:
    provider: "lambdalabs"
    api_key: "test-key"
  invalid: [unclosed
`

	err = os.WriteFile(configPath, []byte(invalidYAML), 0o600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	_ = os.Setenv("HOME", tempDir)

	_, err = LoadConfig()
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("Expected 'failed to parse config file' error, got: %v", err)
	}
}

func TestLoadConfig_DefaultSettings(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".bcloud", "credentials.yaml")

	err := os.MkdirAll(filepath.Dir(configPath), 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `
credentials:
  test-cred:
    provider: "lambdalabs"
    APIKey: "test-key"
    default_location: "us-west-1"
`

	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", originalHome) }()
	_ = os.Setenv("HOME", tempDir)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.Settings.OutputFormat != "yaml" {
		t.Errorf("Expected default output_format 'yaml', got '%s'", config.Settings.OutputFormat)
	}
	if config.Settings.DefaultTimeout != "5m" {
		t.Errorf("Expected default default_timeout '5m', got '%s'", config.Settings.DefaultTimeout)
	}
}

func TestCredentialEntry_DecodeFromMap_MissingProvider(t *testing.T) {
	var entry CredentialEntry
	m := map[string]any{
		"api_key": "test-key",
	}

	err := entry.decodeFromMap(m, "test-key")
	if err == nil {
		t.Fatal("Expected error for missing provider")
	}
	if !strings.Contains(err.Error(), "missing 'provider'") {
		t.Errorf("Expected 'missing provider' error, got: %v", err)
	}
}

func TestCredentialEntry_DecodeFromMap_InvalidProvider(t *testing.T) {
	var entry CredentialEntry
	m := map[string]any{
		"provider": 123,
		"api_key":  "test-key",
	}

	err := entry.decodeFromMap(m, "test-key")
	if err == nil {
		t.Fatal("Expected error for invalid provider")
	}
	if !strings.Contains(err.Error(), "invalid 'provider'") {
		t.Errorf("Expected 'invalid provider' error, got: %v", err)
	}
}

func TestCredentialEntry_DecodeFromMap_UnknownProvider(t *testing.T) {
	var entry CredentialEntry
	m := map[string]any{
		"provider": "unknown-provider",
		"api_key":  "test-key",
	}

	err := entry.decodeFromMap(m, "test-key")
	if err == nil {
		t.Fatal("Expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown provider: unknown-provider") {
		t.Errorf("Expected 'unknown provider' error, got: %v", err)
	}
}

func TestCredentialEntry_DecodeFromMap_KeyAsRefID(t *testing.T) {
	var entry CredentialEntry
	m := map[string]any{
		"provider":         "lambdalabs",
		"APIKey":           "test-key",
		"default_location": "us-west-1",
	}

	err := entry.decodeFromMap(m, "my-lambda-key")
	if err != nil {
		t.Fatalf("decodeFromMap failed: %v", err)
	}

	if entry.Provider != "lambdalabs" {
		t.Errorf("Expected provider 'lambdalabs', got '%s'", entry.Provider)
	}
	if entry.Value.GetReferenceID() != "my-lambda-key" {
		t.Errorf("Expected ref_id 'my-lambda-key', got '%s'", entry.Value.GetReferenceID())
	}

	wrapper, ok := entry.Value.(*LambdaLabsCredentialWrapper)
	if !ok {
		t.Fatal("Expected LambdaLabsCredentialWrapper")
	}
	if wrapper.LambdaLabsCredential.RefID != "my-lambda-key" {
		t.Errorf("Expected embedded RefID 'my-lambda-key', got '%s'", wrapper.LambdaLabsCredential.RefID)
	}
}

func TestCredentialEntry_DecodeFromMap_ExplicitRefID(t *testing.T) {
	var entry CredentialEntry
	m := map[string]any{
		"provider":         "lambdalabs",
		"APIKey":           "test-key",
		"RefID":            "explicit-ref-id",
		"default_location": "us-west-1",
	}

	err := entry.decodeFromMap(m, "yaml-key")
	if err != nil {
		t.Fatalf("decodeFromMap failed: %v", err)
	}

	if entry.Value.GetReferenceID() != "explicit-ref-id" {
		t.Errorf("Expected ref_id 'explicit-ref-id', got '%s'", entry.Value.GetReferenceID())
	}
}

func TestCredentialEntry_JSONMarshalUnmarshal(t *testing.T) {
	original := CredentialEntry{
		Provider: "lambdalabs",
		Value: &LambdaLabsCredentialWrapper{
			LambdaLabsCredential: &lambdalabs.LambdaLabsCredential{
				RefID:  "test-ref",
				APIKey: "test-key",
			},
			DefaultLocation: "us-west-1",
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled CredentialEntry
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if unmarshaled.Provider != original.Provider {
		t.Errorf("Provider mismatch: expected %s, got %s", original.Provider, unmarshaled.Provider)
	}
	if unmarshaled.Value.GetReferenceID() != original.Value.GetReferenceID() {
		t.Errorf("RefID mismatch: expected %s, got %s", original.Value.GetReferenceID(), unmarshaled.Value.GetReferenceID())
	}
}

func TestCredentialEntry_YAMLMarshalUnmarshal(t *testing.T) {
	original := CredentialEntry{
		Provider: "lambdalabs",
		Value: &LambdaLabsCredentialWrapper{
			LambdaLabsCredential: &lambdalabs.LambdaLabsCredential{
				RefID:  "test-ref",
				APIKey: "test-key",
			},
			DefaultLocation: "us-west-1",
		},
	}

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("YAML marshal failed: %v", err)
	}

	var unmarshaled CredentialEntry
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("YAML unmarshal failed: %v", err)
	}

	if unmarshaled.Provider != original.Provider {
		t.Errorf("Provider mismatch: expected %s, got %s", original.Provider, unmarshaled.Provider)
	}
	if unmarshaled.Value.GetReferenceID() != original.Value.GetReferenceID() {
		t.Errorf("RefID mismatch: expected %s, got %s", original.Value.GetReferenceID(), unmarshaled.Value.GetReferenceID())
	}
}

func TestCredentialEntry_EncodeToMap_NilValue(t *testing.T) {
	entry := CredentialEntry{
		Provider: "lambdalabs",
		Value:    nil,
	}

	_, err := entry.encodeToMap()
	if err == nil {
		t.Fatal("Expected error for nil credential value")
	}
	if !strings.Contains(err.Error(), "nil credential Value") {
		t.Errorf("Expected 'nil credential Value' error, got: %v", err)
	}
}

func TestDefaultLocationProvider_Interface(t *testing.T) {
	wrapper := &LambdaLabsCredentialWrapper{
		DefaultLocation: "us-west-1",
	}

	provider, ok := interface{}(wrapper).(DefaultLocationProvider)
	if !ok {
		t.Fatal("LambdaLabsCredentialWrapper should implement DefaultLocationProvider")
	}

	if provider.GetDefaultLocation() != "us-west-1" {
		t.Errorf("Expected default location 'us-west-1', got '%s'", provider.GetDefaultLocation())
	}
}

func TestProviderRegistry_LambdaLabs(t *testing.T) {
	factory, exists := providerRegistry["lambdalabs"]
	if !exists {
		t.Fatal("lambdalabs provider not registered")
	}

	cred := factory()
	if cred == nil {
		t.Fatal("Factory returned nil credential")
	}

	if cred.GetCloudProviderID() != "lambda-labs" {
		t.Errorf("Expected provider ID 'lambda-labs', got '%s'", cred.GetCloudProviderID())
	}

	wrapper, ok := cred.(*LambdaLabsCredentialWrapper)
	if !ok {
		t.Fatal("Expected LambdaLabsCredentialWrapper")
	}

	if wrapper.LambdaLabsCredential == nil {
		t.Fatal("Embedded LambdaLabsCredential is nil")
	}
}

func TestProviderRegistry_Nebius(t *testing.T) {
	factory, exists := providerRegistry["nebius"]
	if !exists {
		t.Fatal("nebius provider not registered")
	}

	cred := factory()
	if cred == nil {
		t.Fatal("Factory returned nil credential")
	}

	if cred.GetCloudProviderID() != "nebius" {
		t.Errorf("Expected provider ID 'nebius', got '%s'", cred.GetCloudProviderID())
	}

	wrapper, ok := cred.(*NebiusCredentialWrapper)
	if !ok {
		t.Fatal("Expected NebiusCredentialWrapper")
	}

	if wrapper.NebiusCredential == nil {
		t.Fatal("Embedded NebiusCredential is nil")
	}
}
