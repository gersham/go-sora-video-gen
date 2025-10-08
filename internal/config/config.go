package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	OpenAIAPIKey string `toml:"openai_api_key"`
	OutputDir    string `toml:"output_dir"`
	Model        string `toml:"model"`
	Duration     string `toml:"duration"`
	Size         string `toml:"size"`
	LastPrompt   string `toml:"last_prompt"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "telemetryos-video-gen.toml"), nil
}

// Load reads the config file from ~/.config/telemetryos-video-gen.toml
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(configPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}

// Save writes the config to ~/.config/telemetryos-video-gen.toml
func Save(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create .config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
