package config

import (
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server server.ServerConfig `yaml:"server"`
	Cache  cache.CacheConfig   `yaml:"cache"`
	// Database DatabaseConfig `yaml:"database"`
	Service ServiceConfig `yaml:"service"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServiceConfig struct {
	Domain string `yaml:"domain" env:"SERVICE_DOMAIN"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL"`
}

func InitConfig() Config {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("read config file")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("unmarshal config file")
	}

	return config
}
