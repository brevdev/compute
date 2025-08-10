package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	lambdalabs "github.com/brevdev/cloud/internal/lambdalabs/v1"
	nebius "github.com/brevdev/cloud/internal/nebius/v1"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"gopkg.in/yaml.v3"
)

const (
	ProviderLambdaLabs = "lambdalabs"
	ProviderNebius     = "nebius"
)

var providerRegistry = map[string]func() v1.CloudCredential{}

func RegisterProvider(id string, factory func() v1.CloudCredential) {
	providerRegistry[id] = factory
}

type LambdaLabsCredentialWrapper struct {
	*lambdalabs.LambdaLabsCredential
	DefaultLocation string `json:"default_location" yaml:"default_location"`
}

func (w *LambdaLabsCredentialWrapper) GetDefaultLocation() string {
	return w.DefaultLocation
}

type NebiusCredentialWrapper struct {
	*nebius.NebiusCredential
	DefaultLocation string `json:"default_location" yaml:"default_location"`
}

func (w *NebiusCredentialWrapper) GetDefaultLocation() string {
	return w.DefaultLocation
}

type DefaultLocationProvider interface {
	GetDefaultLocation() string
}

func init() {
	RegisterProvider(ProviderLambdaLabs, func() v1.CloudCredential {
		return &LambdaLabsCredentialWrapper{
			LambdaLabsCredential: &lambdalabs.LambdaLabsCredential{},
		}
	})
	RegisterProvider(ProviderNebius, func() v1.CloudCredential {
		return &NebiusCredentialWrapper{
			NebiusCredential: &nebius.NebiusCredential{},
		}
	})
}

type CredentialEntry struct {
	Provider string             `json:"provider" yaml:"provider"`
	Value    v1.CloudCredential `json:"-" yaml:"-"`
}

func (c *CredentialEntry) decodeFromMap(m map[string]any, yamlKey string) error {
	rawProv, ok := m["provider"]
	if !ok {
		return fmt.Errorf("missing 'provider'")
	}
	provider, ok := rawProv.(string)
	if !ok || provider == "" {
		return fmt.Errorf("invalid 'provider'")
	}
	factory, ok := providerRegistry[provider]
	if !ok {
		return fmt.Errorf("unknown provider: %s", provider)
	}
	cred := factory()

	if _, hasRefID := m["ref_id"]; !hasRefID {
		m["ref_id"] = yamlKey
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, cred); err != nil {
		return err
	}

	c.Provider = provider
	c.Value = cred
	return nil
}

func (c CredentialEntry) encodeToMap() (map[string]any, error) {
	if c.Value == nil {
		return nil, fmt.Errorf("nil credential Value")
	}
	b, err := json.Marshal(c.Value) // serialize provider-specific fields
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	out["provider"] = string(c.Value.GetCloudProviderID())
	return out, nil
}

func (c *CredentialEntry) UnmarshalJSON(b []byte) error {
	m := map[string]any{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	return c.decodeFromMap(m, "")
}

func (c CredentialEntry) MarshalJSON() ([]byte, error) {
	m, err := c.encodeToMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (c *CredentialEntry) UnmarshalYAML(n *yaml.Node) error {
	m := map[string]any{}
	if err := n.Decode(&m); err != nil {
		return err
	}
	return c.decodeFromMap(m, "")
}

func (c CredentialEntry) MarshalYAML() (interface{}, error) {
	m, err := c.encodeToMap()
	if err != nil {
		return nil, err
	}
	return m, nil // let yaml encode the map
}

type Config struct {
	Credentials map[string]CredentialEntry `json:"credentials" yaml:"credentials"`
	Settings    Settings                   `json:"settings" yaml:"settings"`
}

type Settings struct {
	OutputFormat   string `yaml:"output_format"`
	DefaultTimeout string `yaml:"default_timeout"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".bcloud", "credentials.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		workingDirConfig := "./bcloud.yaml"
		if _, err := os.Stat(workingDirConfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("no configuration found at %s or %s", configPath, workingDirConfig)
		}
		configPath = workingDirConfig
	}

	data, err := os.ReadFile(configPath) // #nosec G304 - configPath is constructed from user home dir and known filename
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var rawConfig struct {
		Credentials map[string]map[string]any `yaml:"credentials"`
		Settings    Settings                  `yaml:"settings"`
	}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	config := Config{
		Credentials: make(map[string]CredentialEntry),
		Settings:    rawConfig.Settings,
	}

	for yamlKey, credData := range rawConfig.Credentials {
		var credEntry CredentialEntry
		if err := credEntry.decodeFromMap(credData, yamlKey); err != nil {
			return nil, fmt.Errorf("failed to parse credential '%s': %w", yamlKey, err)
		}
		config.Credentials[yamlKey] = credEntry
	}

	if config.Settings.OutputFormat == "" {
		config.Settings.OutputFormat = "yaml"
	}
	if config.Settings.DefaultTimeout == "" {
		config.Settings.DefaultTimeout = "5m"
	}

	return &config, nil
}
