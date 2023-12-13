package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256Base58UrlShortener(t *testing.T) {
	links := []struct {
		link     string
		expected string
	}{
		{link: "https://www.google.com", expected: "Cc4PMy5i"},
		{link: "https://www.facebook.com", expected: "FeE4Jcxv"},
		{link: "https://www.twitter.com", expected: "GWWZX9GW"},
		{link: "https://www.instagram.com", expected: "Ah96nhHU"},
		{link: "https://www.linkedin.com", expected: "EGpKkqQM"},
	}

	shortener := &Sha256Base58UrlShortener{}

	for _, link := range links {
		shortLink, err := shortener.Shorten(link.link)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, shortLink, link.expected)
	}
}
