package handlers

import (
	//"fmt"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/data"
)

type GetHandler struct {
	shortLinkRepo data.ShortLinkRepo
}

func NewGetHandler(repo data.ShortLinkRepo) *GetHandler {
	handler := new(GetHandler)
	handler.shortLinkRepo = repo
	return handler
}

func (handler *GetHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	id := req.PathValue("id")
	pathegmentCount := len(strings.Split(req.URL.Path, "/"))

	if id == "" || pathegmentCount <= 1 || pathegmentCount > 2 {
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	url := handler.shortLinkRepo.GetURL(id)

	if url == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
