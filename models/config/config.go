package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
	"github.com/ChristianLapinig/aem-local-cli/models/environment"
)

const (
	configFileName = "config.json"
)

type Config struct {
	Environments []environment.Environment `json:"environments"`
}

func CreateConfigFile(path string) error {
	config := Config{
		Environments: []environment.Environment{},
	}
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	configFile := filepath.Join(path, configFileName)
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func GetTempFolderPath() (string, error) {
	data, err := utils.LoadMarkerFile()
	path := filepath.Join(strings.TrimSpace(string(data)), "temp")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}

func GetConfigPath() (string, error) {
	data, err := utils.LoadMarkerFile()
	path := filepath.Join(strings.TrimSpace(string(data)), "config.json")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return &Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}

	var config *Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{}, err
	}
	return config, err
}

func UpdateConfig(path string, config *Config) error {
	out, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
