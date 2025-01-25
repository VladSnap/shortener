package handlers

import (
	//"fmt"
	"io"
	"net/http"
	"strings"

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
		http.Error(res, "Http method not POST", http.StatusBadRequest)
		return
	}

	if req.URL.Path != "/" {
		http.Error(res, "Incorrect request path", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get("content-type")

	if ct != "text/plain" && ct != "text/plain; charset=utf-8" && ct != "application/x-gzip" {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	fullURL := string(body)

	if err != nil || fullURL == "" {
		http.Error(res, "Required url", http.StatusBadRequest)
		return
	}

	fullURL = strings.TrimSuffix(fullURL, "\r")
	fullURL = strings.TrimSuffix(fullURL, "\n")
	verifyRes, urlIsValid := urlverifier.NewVerifier().Verify(fullURL)

	if urlIsValid != nil || !verifyRes.IsURL || !verifyRes.IsRFC3986URL {
		http.Error(res, "Full URL verify error", http.StatusBadRequest)
		return
	}

	shortLink, err := handler.service.CreateShortLink(fullURL)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(handler.baseURL + "/" + shortLink))
}
