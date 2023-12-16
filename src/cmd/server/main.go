package main

import (
	v1 "github.com/lilydoar/short-n-sweet/src/internal/api/v1"
	"github.com/lilydoar/short-n-sweet/src/internal/config"
	"github.com/lilydoar/short-n-sweet/src/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogging(cfg config.LoggingConfig) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}

	zerolog.SetGlobalLevel(level)

	ctxLogger := log.With().Logger()
	zerolog.DefaultContextLogger = &ctxLogger
}

func main() {
	cfg := config.InitConfig()
	initLogging(cfg.Logging)

	router := v1.InitRouter(cfg)
	server.ListenAndServeWithGracefulShutdown(router, cfg.Server)
}
