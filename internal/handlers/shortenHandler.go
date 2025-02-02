package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	urlverifier "github.com/davidmytton/url-verifier"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenHandler struct {
	service services.ShorterService
	baseURL string
}

func NewShortenHandler(service services.ShorterService, baseURL string) *ShortenHandler {
	handler := new(ShortenHandler)
	handler.service = service
	handler.baseURL = baseURL
	return handler
}

func (handler *ShortenHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Http method not POST", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get("content-type")

	if ct != "application/json" && ct != "application/x-gzip" && ct != "application/json; charset=utf-8" {
		http.Error(res, "Incorrect content-type:"+ct, http.StatusBadRequest)
		return
	}

	var request ShortenRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		http.Error(res, "Required url", http.StatusBadRequest)
		return
	}

	request.URL = strings.TrimSuffix(request.URL, "\r")
	request.URL = strings.TrimSuffix(request.URL, "\n")
	verifyRes, urlIsValid := urlverifier.NewVerifier().Verify(request.URL)

	if urlIsValid != nil || !verifyRes.IsURL || !verifyRes.IsRFC3986URL {
		http.Error(res, "Full URL verify error", http.StatusBadRequest)
		return
	}

	shortLink, err := handler.service.CreateShortLink(request.URL)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := ShortenResponse{Result: handler.baseURL + "/" + shortLink}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(result)

	if err != nil {
		log.Zap.Errorf("failed write to response: %w", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
