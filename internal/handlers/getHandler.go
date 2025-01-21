package handlers

import (
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/services"
)

type GetHandler struct {
	service services.ShorterService
}

func NewGetHandler(service services.ShorterService) *GetHandler {
	handler := new(GetHandler)
	handler.service = service
	return handler
}

func (handler *GetHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortID := req.PathValue("id")
	pathegmentCount := len(strings.Split(req.URL.Path, "/"))

	if shortID == "" || pathegmentCount <= 1 || pathegmentCount > 2 {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	url := handler.service.GetURL(shortID)

	if url == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
