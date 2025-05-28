package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

// LoadConfig loads configuration from file
func LoadConfig(configFile string) error {
	// If the configuration file does not exist, create a default configuration
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig(configFile)
		}
		return fmt.Errorf("Failed to read configuration file: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, Global); err != nil {
		return fmt.Errorf("Failed to parse configuration file: %w", err)
	}
	log.Info().Msg("Configuration file loaded successfully")
	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configFile string) error {
	log.Info().Msg("Configuration file does not exist, creating default configuration file")

	if err := os.WriteFile(configFile, []byte(DefaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("Failed to write default configuration file: %w", err)
	}

	log.Info().Msg("Default configuration file created successfully")
	log.Info().Msg(fmt.Sprintf("Please edit the configuration file: %s", configFile))
	os.Exit(0)
	return nil
}
