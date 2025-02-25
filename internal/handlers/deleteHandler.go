package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/validation"
)

type DeleteHandler struct {
	service ShorterService
	mtx     sync.Mutex
}

func NewDeleteHandler(service ShorterService) *DeleteHandler {
	handler := new(DeleteHandler)
	handler.service = service
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

	handler.mtx.Lock()
	err := handler.service.DeleteBatch(req.Context(), shortURLs, userID)
	handler.mtx.Unlock()

	if err != nil {
		http.Error(res, fmt.Errorf("failed DeleteBatch: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	res.WriteHeader(http.StatusAccepted)
}
