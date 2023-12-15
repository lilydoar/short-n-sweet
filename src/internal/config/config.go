package config

import (
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server server.ServerConfig `json:"server" yaml:"server" toml:"server"`
	Cache  cache.CacheConfig   `json:"cache" yaml:"cache" toml:"cache"`
	// Database DatabaseConfig `json:"database" yaml:"database" toml:"database"`
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
