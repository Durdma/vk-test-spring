package config

import (
	"github.com/spf13/viper"
	filepath2 "path/filepath"
	"time"
)

const (
	defaultHttpPort      = "8080"
	defaultHttpRWTimeout = 10 * time.Second
	defaultLoggerLevel   = 5
)

type Config struct {
	// PostgreSQL  PostgreSQLConfig
	HTTP        HTTPConfig
	LoggerLevel int
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgreSQLConfig struct {
}

func Init(path string) (*Config, error) {
	setDefaults()

	var cfg Config

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("logger.level", &cfg.LoggerLevel); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(filepath string) error {
	path := filepath2.Dir(filepath)
	name := filepath2.Base(filepath)

	viper.AddConfigPath(path)
	viper.SetConfigName(name)

	return viper.ReadInConfig()
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("logger.level", defaultLoggerLevel)
}
