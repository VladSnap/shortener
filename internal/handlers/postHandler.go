package handlers

import (
	//"fmt"
	"io"
	"net/http"

	"github.com/VladSnap/shortener/internal/data"
)

type PostHandler struct {
	shortLinkRepo data.ShortLinkRepo
	baseURL       string
}

func NewPostHandler(repo data.ShortLinkRepo, baseURL string) *PostHandler {
	handler := new(PostHandler)
	handler.shortLinkRepo = repo
	handler.baseURL = baseURL
	return handler
}

func (handler *PostHandler) Handle(res http.ResponseWriter, req *http.Request) {
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

	if err != nil || string(body) == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortLink := handler.shortLinkRepo.CreateShortLink(string(body))

	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(handler.baseURL + shortLink))
}
