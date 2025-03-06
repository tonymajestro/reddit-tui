package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"reddittui/utils"

	"github.com/BurntSushi/toml"
)

const (
	configFilename    = "reddittui.toml"
	defaultDomainName = "old.reddit.com"
	defaultServerType = "old"
)

type Config struct {
	Core   CoreConfig   `toml:"core"`
	Filter FilterConfig `toml:"filter"`
	Client ClientConfig `toml:"client"`
	Server ServerConfig `toml:"server"`
}

type CoreConfig struct {
	BypassCache   bool
	LogLevel      string
	ClientTimeout int // Legacy
}

type FilterConfig struct {
	Keywords   []string
	Subreddits []string
}

type ClientConfig struct {
	TimeoutSeconds  int
	CacheTtlSeconds int
}

type ServerConfig struct {
	Domain string
	Type   string
}

func NewConfig() Config {
	return Config{
		Core: CoreConfig{
			BypassCache:   false,
			LogLevel:      "Warn",
			ClientTimeout: 10,
		},
		Server: ServerConfig{
			Domain: defaultDomainName,
			Type:   defaultServerType,
		},
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
	if meta.IsDefined("core", "bypassCache") {
		left.Core.BypassCache = right.Core.BypassCache
	}

	if meta.IsDefined("core", "logLevel") {
		left.Core.LogLevel = right.Core.LogLevel
	}

	if meta.IsDefined("core", "clientTimeout") {
		left.Core.ClientTimeout = right.Core.ClientTimeout
	}

	if meta.IsDefined("filter", "keywords") {
		left.Filter.Keywords = right.Filter.Keywords
	}

	if meta.IsDefined("filter", "subreddits") {
		left.Filter.Subreddits = right.Filter.Subreddits
	}

	if meta.IsDefined("client", "timeoutSeconds") {
		left.Client.TimeoutSeconds = right.Client.TimeoutSeconds
	}

	if meta.IsDefined("client", "cacheTtlSeconds") {
		left.Client.CacheTtlSeconds = right.Client.CacheTtlSeconds
	}

	if meta.IsDefined("server", "domain") {
		left.Server.Domain = right.Server.Domain
	}

	if meta.IsDefined("server", "type") {
		left.Server.Type = right.Server.Type
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
