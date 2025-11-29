package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	ProfilesDir string `mapstructure:"profiles_dir"`
	Profile     string `mapstructure:"profile"`
	UserID      string `mapstructure:"user_id"`
	CharacterID string `mapstructure:"character_id"`
}

// LoadConfig loads configuration from file or returns default config
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Set config file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add current directory as config path
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("profiles_dir", "")
	viper.SetDefault("profile", "")
	viper.SetDefault("user_id", "")
	viper.SetDefault("character_id", "")

	// Read config file (ignore error if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	viper.Set("profiles_dir", cfg.ProfilesDir)
	viper.Set("profile", cfg.Profile)
	viper.Set("user_id", cfg.UserID)
	viper.Set("character_id", cfg.CharacterID)

	// Set config file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Ensure config directory exists
	configPath := "config.yaml"
	configDir := filepath.Dir(configPath)
	if configDir != "." && configDir != "" {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		// If file doesn't exist, try SafeWriteConfigAs
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
	}

	return nil
}
