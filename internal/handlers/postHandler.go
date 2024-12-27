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

	ct := req.Header.Get("content-type")

	if ct != "text/plain" && ct != "text/plain; charset=utf-8" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortLink := data.CreateShortLink(string(body))

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Write([]byte("http://localhost:8080/" + shortLink))
}
