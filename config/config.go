package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"reddittui/utils"

	"github.com/BurntSushi/toml"
)

const configFilename = "reddit-tui.toml"

const defaultConfiguration = `
#
# Default configuration for reddit-tui.
# Uncomment to configure
#

# bypassCache = false
# logLevel = Info
`

type Config struct {
	BypassCache bool   `toml:"bypassCache"`
	LogLevel    string `toml:"logLevel"`
}

func NewConfig() Config {
	return Config{
		BypassCache: false,
		LogLevel:    "Info",
	}
}

func LoadConfig() (Config, error) {
	defaultConfig := NewConfig()

	configDir, err := utils.GetConfigDir()
	if err != nil {
		slog.Warn("Could not get config directory", "error", err)
		return defaultConfig, err
	}

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		slog.Warn("Could not make config directory", "error", err)
		return defaultConfig, err
	}

	configPath := filepath.Join(configDir, configFilename)
	configFile, err := os.Open(configPath)
	if os.IsNotExist(err) {
		createConfigFile(configPath)
		return defaultConfig, err
	} else if err != nil {
		slog.Warn("Could not open config file", "error", err)
		return defaultConfig, err
	}

	defer configFile.Close()

	var configFromFile Config
	decoder := toml.NewDecoder(configFile)
	meta, err := decoder.Decode(&configFromFile)
	if err != nil {
		slog.Warn("Could not decode config file", "error", err)
		return defaultConfig, err
	}

	mergedConfig := mergeConfig(defaultConfig, configFromFile, meta)
	return mergedConfig, err
}

// Merge right config into left
func mergeConfig(left, right Config, meta toml.MetaData) Config {
	if meta.IsDefined("bypassCache") {
		left.BypassCache = right.BypassCache
	}

	if meta.IsDefined("logLevel") {
		left.LogLevel = right.LogLevel
	}

	return left
}

func createConfigFile(configFilePath string) error {
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	_, err = configFile.WriteString(defaultConfiguration)
	return err
}
