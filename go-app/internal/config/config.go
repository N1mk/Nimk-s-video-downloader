package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type RawConfig struct {
	DownloadPath   string `json:"download_path"`
	MaxRetryCount  string `json:"max_retry_count"`
	AddToAutostart string `json:"add_to_autostart"`
}

type Config struct {
	DownloadPath   string `json:"download_path"`
	MaxRetryCount  int    `json:"max_retry_count"`
	AddToAutostart bool   `json:"add_to_autostart"`
}

const (
	ConfigFileName        string = "config.json"
	DefaultDownloadPath   string = ""
	DefaultMaxRetryCount  int    = 2
	DefaultAddToAutostart        = true
)

type ConfigReader struct {
	configPath string
}

func NewConfigReader(exePath string) (configReader *ConfigReader, err error) {
	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, ConfigFileName)

	if _, err := os.Stat(configPath); err != nil {
		os.Create(configPath)

		config := &Config{DownloadPath: DefaultDownloadPath, MaxRetryCount: DefaultMaxRetryCount, AddToAutostart: DefaultAddToAutostart}

		data, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("config marshalling error: %w", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("writing config file error: %w", err)
		}
	}

	return &ConfigReader{configPath: configPath}, nil
}

func (cr *ConfigReader) GetConfig() (*Config, error) {
	data, err := os.ReadFile(cr.configPath)
	if err != nil {
		return nil, fmt.Errorf("config file reading error: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config unmarshalling error: %w", err)
	}

	return &config, nil
}

func (cr *ConfigReader) SetConfig(config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("config marshalling error: %w", err)
	}

	if err := os.WriteFile(cr.configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config file error: %w", err)
	}

	return nil
}

func (cr *ConfigReader) GetConfigJSON() ([]byte, error) {
	data, err := os.ReadFile(cr.configPath)
	if err != nil {
		return nil, fmt.Errorf("config file reading error: %w", err)
	}

	return data, nil
}

func (cr *ConfigReader) SetConfigJSON(data []byte) error {
	if err := os.WriteFile(cr.configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config file error: %w", err)
	}

	return nil
}
