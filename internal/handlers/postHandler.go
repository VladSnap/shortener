package handlers

import (
	//"fmt"
	"io"
	"net/http"

	"github.com/VladSnap/shortener/internal/services"
	urlverifier "github.com/davidmytton/url-verifier"
)

type PostHandler struct {
	service services.ShorterService
	baseURL string
}

func NewPostHandler(service services.ShorterService, baseURL string) *PostHandler {
	handler := new(PostHandler)
	handler.service = service
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
	fullUrl := string(body)

	if err != nil || fullUrl == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	verifyRes, urlIsValid := urlverifier.NewVerifier().Verify(fullUrl)

	if urlIsValid != nil || !verifyRes.IsURL || !verifyRes.IsRFC3986URL {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortLink, err := handler.service.CreateShortLink(fullUrl)

	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(handler.baseURL + "/" + shortLink))
}
