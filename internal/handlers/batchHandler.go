package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/VladSnap/shortener/internal/validation"
)

type ShortenRowRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenRowResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchHandler struct {
	service ShorterService
	baseURL string
}

func NewBatchHandler(service ShorterService, baseURL string) *BatchHandler {
	handler := new(BatchHandler)
	handler.service = service
	handler.baseURL = baseURL
	return handler
}

func (handler *BatchHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Http method not POST", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get("content-type")

	if !strings.Contains(ct, HeaderApplicationJSON) && !strings.Contains(ct, HeaderApplicationXgzip) {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	var requestRows []ShortenRowRequest

	if err := json.NewDecoder(req.Body).Decode(&requestRows); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	links := make([]*services.OriginalLink, 0, len(requestRows))
	for _, r := range requestRows {
		r.OriginalURL = strings.TrimSuffix(r.OriginalURL, "\r")
		r.OriginalURL = strings.TrimSuffix(r.OriginalURL, "\n")
		if err := validation.ValidateURL(r.OriginalURL, "OriginalURL"); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		lin := &services.OriginalLink{
			CorelationID: r.CorrelationID,
			URL:          r.OriginalURL,
		}
		links = append(links, lin)
	}

	if len(requestRows) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	shortedLinks, err := handler.service.CreateShortLinkBatch(req.Context(), links)

	if err != nil {
		http.Error(res, fmt.Errorf("failed CreateShortLinkBatch: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	responseRows := make([]*ShortenRowResponse, 0, len(shortedLinks))
	for _, sl := range shortedLinks {
		rr := &ShortenRowResponse{
			CorrelationID: sl.CorelationID,
			ShortURL:      handler.baseURL + "/" + sl.URL,
		}

		responseRows = append(responseRows, rr)
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(responseRows)

	if err != nil {
		log.Zap.Errorf(ErrFailedWriteToResponse, err)
		return
	}
}
