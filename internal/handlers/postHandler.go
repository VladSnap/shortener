package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/validation"
)

type PostHandler struct {
	service ShorterService
	baseURL string
}

func NewPostHandler(service ShorterService, baseURL string) *PostHandler {
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
	if !strings.Contains(ct, "text/plain") && !strings.Contains(ct, HeaderApplicationXgzip) {
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
	if err := validation.ValidateURL(fullURL, "req.Body"); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID := "d1a8485a-430a-49f4-92ba-50886e1b07c6"
	shortLink, err := handler.service.CreateShortLink(req.Context(), fullURL, userID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "text/plain")
	if shortLink.IsDuplicated {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	_, err = res.Write([]byte(handler.baseURL + "/" + shortLink.URL))

	if err != nil {
		log.Zap.Errorf(ErrFailedWriteToResponse, err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
