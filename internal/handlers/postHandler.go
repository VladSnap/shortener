package handlers

import (
	//"fmt"
	"io"
	"net/http"

	"github.com/VladSnap/shortener/internal/data"
)

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.URL.Path != "/" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Header.Get("content-type") != "text/plain" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortLink := data.CreateShortLink(string(body))

	res.Header().Set("content-type", "text/plain")
	res.Write([]byte("http://localhost:8080/" + shortLink))
}
