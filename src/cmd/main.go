package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/config"
	"github.com/lilydoar/short-n-sweet/src/internal/handlers"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	config := config.InitConfig()

	redisCache := cache.InitRedisCache(config.Cache)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Short-n-Sweet URL shortener"))
	})
	router.Handle("/encode", handlers.EncodeHandler(redisCache, config)).Methods("POST")
	router.Handle("/decode/{encodedUrl}", handlers.DecodeHandler(redisCache)).Methods("GET")

	addr := config.Server.Host + ":" + config.Server.Port

	log.Info("Starting server at " + addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
