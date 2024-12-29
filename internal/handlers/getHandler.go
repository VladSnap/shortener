package handlers

import (
	//"fmt"
	"net/http"

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

	if id == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	url := handler.shortLinkRepo.GetURL(id)

	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
