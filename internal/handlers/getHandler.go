package handlers

import (
	//"fmt"
	"net/http"

	"github.com/VladSnap/shortener/internal/data"
)

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	id := req.PathValue("id")

	if id == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

	url := data.GetURL(id)

	res.Header().Set("Location", url)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
