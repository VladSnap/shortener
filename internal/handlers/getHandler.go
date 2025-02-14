package handlers

import (
	"net/http"
	"strings"
)

type GetHandler struct {
	service ShorterService
}

func NewGetHandler(service ShorterService) *GetHandler {
	handler := new(GetHandler)
	handler.service = service
	return handler
}

func (handler *GetHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Http method not GET", http.StatusBadRequest)
		return
	}

	shortID := req.PathValue("id")
	if !validatePath(req.URL.Path) || shortID == "" {
		http.Error(res, "Request path incorrect", http.StatusBadRequest)
		return
	}

	url, err := handler.service.GetURL(shortID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if url == "" {
		http.Error(res, "Url not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func validatePath(path string) bool {
	segments := strings.Split(path, "/")
	return len(segments) == 2 && segments[1] != ""
}
