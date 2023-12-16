package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/config"
	"github.com/lilydoar/short-n-sweet/src/internal/shortener"
	"github.com/rs/zerolog"
)

type EncodeRequest struct {
	Url string `json:"url"`
}

type EncodeResponse struct {
	Url string `json:"url"`
}

func EncodeHandler(cache cache.Cache, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())

		var request EncodeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			l.Error().Err(err).Msg("invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request body"))
			return
		}

		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_url", request.Url)
		})

		shortener := shortener.Sha256Base58UrlShortener{}
		encodedUrl, err := shortener.Shorten(request.Url)
		if err != nil {
			l.Error().Err(err).Msg("encoding url")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("encoded_url", encodedUrl)
		})

		l.Debug().Msg("encoded url")

		cache.Set(encodedUrl, request.Url)

		l.Debug().Msg("cached url")

		responseData, err := json.Marshal(EncodeResponse{Url: cfg.Service.Domain + "/decode/" + encodedUrl})
		if err != nil {
			l.Error().Err(err).Msg("marshalling json response")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	})
}

func DecodeHandler(cache cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		encodedUrl := vars["encodedUrl"]

		l := zerolog.Ctx(r.Context())
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("encoded_url", encodedUrl)
		})

		decodedUrl, err := cache.Get(encodedUrl)
		if err != nil {
			l.Error().Err(err).Msg("retrieving url from cache")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		l.Debug().Str("decoded_url", decodedUrl).Msg("decoded url")

		http.Redirect(w, r, decodedUrl, http.StatusMovedPermanently)
	})
}
