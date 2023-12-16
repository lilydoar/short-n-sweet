package v1

import (
	"github.com/gorilla/mux"
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/config"
	"github.com/lilydoar/short-n-sweet/src/internal/handlers"
	"github.com/lilydoar/short-n-sweet/src/internal/middleware"
)

func InitRouter(cfg config.Config) *mux.Router {
	cache := cache.InitRedisCache(cfg.Cache)

	router := mux.NewRouter()

	router.Use(middleware.LoggingContextMiddleware)
	router.Use(middleware.RequestIdMiddleWare)
	router.Use(middleware.RequestLoggingMiddleware)

	router.Methods("POST").Path("/encode").Handler(handlers.EncodeHandler(cache, cfg))
	router.Methods("GET").Path("/decode/{encodedUrl}").Handler(handlers.DecodeHandler(cache))

	return router
}
