package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/validation"
)

type DeleteHandler struct {
	service ShorterService
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

	go func() {
		err := handler.service.DeleteBatch(req.Context(), shortURLs, userID)
		if err != nil {
			log.Zap.Errorf("failed DeleteBatch: %w", err)
			return
		}
	}()

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	res.WriteHeader(http.StatusAccepted)
}
