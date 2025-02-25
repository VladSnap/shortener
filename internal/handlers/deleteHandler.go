package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/VladSnap/shortener/internal/validation"
)

type DeleteHandler struct {
	deleteWorker DeleterWorker
}

type DeleterWorker interface {
	Close() error
	AddToDelete(shortIDs chan services.DeleteShortID)
	RunWork()
}

func NewDeleteHandler(deleteWorker DeleterWorker) *DeleteHandler {
	handler := new(DeleteHandler)
	handler.deleteWorker = deleteWorker
	return handler
}

func (handler *DeleteHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Http method not DELETE", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(HeaderContentType)
	if !strings.Contains(ct, HeaderApplicationJSONValue) && !strings.Contains(ct, HeaderApplicationXgzipValue) {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	var shortURLs []string
	if err := json.NewDecoder(req.Body).Decode(&shortURLs); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	for _, surl := range shortURLs {
		if err := validation.ValidateShortURL(surl); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if len(shortURLs) == 0 {
		res.WriteHeader(http.StatusNotAcceptable)
		return
	}

	userID := ""
	if value, ok := req.Context().Value(constants.UserIDContextKey).(string); ok {
		userID = value
	}

	go func() {
		toDeleteChan := make(chan services.DeleteShortID, 100)
		defer close(toDeleteChan)

		for _, url := range shortURLs {
			deleteSID := services.NewDeleteShortID(url, userID)
			toDeleteChan <- deleteSID
		}

		handler.deleteWorker.AddToDelete(toDeleteChan)
	}()

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	res.WriteHeader(http.StatusAccepted)
}
