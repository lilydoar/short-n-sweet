package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
)

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST"`
	Port string `yaml:"port" env:"SERVER_PORT"`
}

func ListenAndServeWithGracefulShutdown(router *mux.Router, cfg ServerConfig) {
	addr := cfg.Host + ":" + cfg.Port

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	log.Info().Str("address", addr).Msg("starting server")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("listen and serve")
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("shutdown server")
	}
}