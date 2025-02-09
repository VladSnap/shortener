package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	urlverifier "github.com/davidmytton/url-verifier"
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

	if ct != "application/json" && ct != "application/x-gzip" && ct != "application/json; charset=utf-8" {
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
		if err := validateRequestRow(r); err != nil {
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

	shortedLinks, err := handler.service.CreateShortLinkBatch(links)

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

func validateRequestRow(rqRow ShortenRowRequest) error {
	if rqRow.OriginalURL == "" {
		return errors.New("required url")
	}

	rqRow.OriginalURL = strings.TrimSuffix(rqRow.OriginalURL, "\r")
	rqRow.OriginalURL = strings.TrimSuffix(rqRow.OriginalURL, "\n")
	verifyRes, urlIsValid := urlverifier.NewVerifier().Verify(rqRow.OriginalURL)

	if urlIsValid != nil || !verifyRes.IsURL || !verifyRes.IsRFC3986URL {
		return errors.New("originalURL verify error")
	}

	return nil
}
