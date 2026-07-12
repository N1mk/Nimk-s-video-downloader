package config_reader

import (
	"encoding/json"
	"fmt"
	"nvd/internal/models"
	"os"
	"path/filepath"
)

type ConfigReader struct {
	configPath string
}

func NewConfigReader(exePath string) *ConfigReader {
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")

	return &ConfigReader{configPath: configPath}
}

func (cr *ConfigReader) GetConfig() (*models.Config, error) {
	data, err := os.ReadFile(cr.configPath)
	if err != nil {
		return nil, fmt.Errorf("config file reading error: %w", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config unmarshalling error: %w", err)
	}

	return &config, nil
}

func (cr *ConfigReader) SetConfig(config *models.Config) error {
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
