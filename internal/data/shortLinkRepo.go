package data

import (
	"github.com/VladSnap/shortener/internal/helpers"
)

const linkKeyLength = 8

var links map[string]string = make(map[string]string)

func CreateShortLink(url string) string {
	key := helpers.RandStringRunes(linkKeyLength)
	links[key] = url
	return key
}

func GetUrl(key string) string {
	return links[key]
}
