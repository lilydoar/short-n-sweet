package shortener

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

type UrlShortener interface {
	Shorten(longUrl string) (string, error)
}

type Sha256Base58UrlShortener struct{}

func (shortener *Sha256Base58UrlShortener) Shorten(longUrl string) (string, error) {
	hash := sha256.Sum256([]byte(longUrl))
	encoded := base58.Encode(hash[:])

	if len(encoded) > 8 {
		return encoded[:8], nil
	} else {
		return encoded, nil
	}
}
