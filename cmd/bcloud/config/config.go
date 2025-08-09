package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Credentials map[string]CredentialConfig `yaml:"credentials"`
	Settings    Settings                    `yaml:"settings"`
}

type CredentialConfig struct {
	Provider          string `yaml:"provider"`
	RefID             string `yaml:"ref_id,omitempty"`
	APIKey            string `yaml:"api_key,omitempty"`
	ServiceAccountKey string `yaml:"service_account_key,omitempty"`
	ProjectID         string `yaml:"project_id,omitempty"`
	DefaultLocation   string `yaml:"default_location,omitempty"`
}

type Settings struct {
	OutputFormat   string `yaml:"output_format"`
	DefaultTimeout string `yaml:"default_timeout"`
}

func (c *CredentialConfig) GetRefID(yamlKey string) string {
	if c.RefID != "" {
		return c.RefID
	}
	return yamlKey
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

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	if config.Settings.OutputFormat == "" {
		config.Settings.OutputFormat = "yaml"
	}
	if config.Settings.DefaultTimeout == "" {
		config.Settings.DefaultTimeout = "5m"
	}

	return &config, nil
}
