package config

import (
	"github.com/spf13/viper"
	filepath2 "path/filepath"
	"runtime"
	"time"
)

const (
	defaultHttpPort      = "8080"
	defaultHttpRWTimeout = 10 * time.Second
	defaultLoggerLevel   = 5
)

type Config struct {
	PostgreSQL PostgreSQLConfig
	HTTP       HTTPConfig
	Logger     LoggerConfig
}

type LoggerConfig struct {
	LoggerLevel    int
	LoggerFilePath string
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgreSQLConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifetime time.Duration
	DriverName            string
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
	if err := viper.UnmarshalKey("logger.level", &cfg.Logger.LoggerLevel); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("loggger.file", &cfg.Logger.LoggerFilePath); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgresql", &cfg.PostgreSQL); err != nil {
		return err
	}

	return nil
}

// TODO Refactor configPath
func parseConfigFile(filepath string) error {

	//dir := filepath2.Dir(filepath)
	//
	//viper.AddConfigPath(dir)
	//
	//name := filepath2.Base(filepath)
	//viper.SetConfigName(name)
	if os := runtime.GOOS; os == "linux" {
		viper.SetConfigFile("/app/configs/main.yaml")
	} else {
		path := filepath2.Dir(filepath)
		name := filepath2.Base(filepath)

		viper.AddConfigPath(path)
		viper.SetConfigName(name)
	}

	return viper.ReadInConfig()
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("logger.level", defaultLoggerLevel)
}
