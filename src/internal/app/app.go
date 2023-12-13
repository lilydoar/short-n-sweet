package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lilydoar/short-n-sweet/src/internal/cache"
	"github.com/lilydoar/short-n-sweet/src/internal/shortener"
	log "github.com/sirupsen/logrus"
)

const host = "http://localhost:8080/"

type App struct {
	CacheService cache.Cache
}

type ShortUrlCreationRequest struct {
	LongUrl string `json:"longUrl"` 
}

type ShortUrlCreationResponse struct {
	ShortUrl string `json:"shortUrl"`
}

func (a *App)CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	var request ShortUrlCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err) 
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	log.WithFields(log.Fields{"method": r.Method, "longUrl": request.LongUrl}).Info("Received request to shorten URL")

	shortener := shortener.Sha256Base58UrlShortener{}
	shortUrl, err := shortener.Shorten(request.LongUrl)
	if err != nil {
		log.WithFields(log.Fields{"longUrl": request.LongUrl}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	a.CacheService.Set(shortUrl, request.LongUrl)

	jsonData, err := json.Marshal(ShortUrlCreationResponse{ShortUrl: host + shortUrl})
	if err != nil {
		log.WithFields(log.Fields{"longUrl": request.LongUrl, "shortUrl": shortUrl}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	log.WithFields(log.Fields{"longUrl": request.LongUrl, "shortUrl": shortUrl}).Info("Shortened URL")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (a *App)HandleShortUrlRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]

	log.WithFields(log.Fields{"method": r.Method, "shortUrl": shortUrl}).Info("Received request to redirect short URL")	

	longUrl, err := a.CacheService.Get(shortUrl)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	log.WithFields(log.Fields{"shortUrl": shortUrl, "longUrl": longUrl}).Info("Redirecting short URL")

	http.Redirect(w, r, longUrl, http.StatusFound)
}