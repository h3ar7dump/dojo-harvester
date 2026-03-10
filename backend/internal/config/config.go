package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Logger    LoggerConfig    `mapstructure:"logger"`
	Scripts   ScriptsConfig   `mapstructure:"scripts"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Platform  PlatformConfig  `mapstructure:"platform"`
}

type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Mode         string `mapstructure:"mode"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

type StorageConfig struct {
	BadgerPath     string `mapstructure:"badger_path"`
	DatasetPath    string `mapstructure:"dataset_path"`
	MinFreeSpaceMB int64  `mapstructure:"min_free_space_mb"`
}

type LoggerConfig struct {
	Level       string `mapstructure:"level"`
	Development bool   `mapstructure:"development"`
}

type ScriptsConfig struct {
	Record      string `mapstructure:"record"`
	Convert     string `mapstructure:"convert"`
	UploadLocal string `mapstructure:"upload_local"`
	TimeoutSecs int    `mapstructure:"timeout_secs"`
}

type TelemetryConfig struct {
	MaxConcurrentConnections int `mapstructure:"max_concurrent_connections"`
	MaxMessageSizeBytes      int `mapstructure:"max_message_size_bytes"`
}

type PlatformConfig struct {
	URL        string `mapstructure:"url"`
	RetryCount int    `mapstructure:"retry_count"`
}

func LoadConfig(path string) (*Config, error) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath("config")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("DOJO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Default configurations
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && path != "" {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "development")
	viper.SetDefault("server.allow_origins", "*")

	viper.SetDefault("storage.badger_path", "./data/badger")
	viper.SetDefault("storage.dataset_path", "./data/datasets")
	viper.SetDefault("storage.min_free_space_mb", 1024)

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.development", true)

	viper.SetDefault("scripts.record", "../toolkit/scripts/record.sh")
	viper.SetDefault("scripts.convert", "../toolkit/scripts/convert.sh")
	viper.SetDefault("scripts.upload_local", "../toolkit/scripts/upload_local.sh")
	viper.SetDefault("scripts.timeout_secs", 3600)

	viper.SetDefault("telemetry.max_concurrent_connections", 50)
	viper.SetDefault("telemetry.max_message_size_bytes", 1048576)

	viper.SetDefault("platform.url", "http://localhost:3000")
	viper.SetDefault("platform.retry_count", 5)
}
