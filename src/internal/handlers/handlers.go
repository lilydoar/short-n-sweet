package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/config"
	"github.com/lilydoar/short-n-sweet/src/internal/shortener"
	"github.com/rs/zerolog/log"
)

type EncodeRequest struct {
	Url string `json:"url"`
}

type EncodeResponse struct {
	Url string `json:"url"`
}

func EncodeHandler(cache cache.Cache, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request EncodeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Error().Err(err).Msg("invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request body"))
			return
		}

		log.Info().Str("method", r.Method).Str("url", request.Url).Msg("received request to encode url")

		shortener := shortener.Sha256Base58UrlShortener{}
		encodedUrl, err := shortener.Shorten(request.Url)
		if err != nil {
			log.Error().Err(err).Str("url", request.Url).Msg("error encoding url")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		cache.Set(encodedUrl, request.Url)

		responseData, err := json.Marshal(EncodeResponse{Url: cfg.Service.Domain + "/decode/" + encodedUrl})
		if err != nil {
			log.Error().Err(err).Msg("error marshalling json response")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		log.Info().Str("method", r.Method).Str("url", request.Url).Str("encodedUrl", encodedUrl).Msg("successfully encoded url")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	})
}

func DecodeHandler(cache cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		encodedUrl := vars["encodedUrl"]

		log.Info().Str("method", r.Method).Str("encodedUrl", encodedUrl).Msg("received request to decode url")

		url, err := cache.Get(encodedUrl)
		if err != nil {
			log.Error().Err(err).Str("encodedUrl", encodedUrl).Msg("error retrieving url from cache")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		log.Info().Str("method", r.Method).Str("encodedUrl", encodedUrl).Str("url", url).Msg("redirecting to url")

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})
}
