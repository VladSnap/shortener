package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/validation"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenHandler struct {
	service ShorterService
	baseURL string
}

func NewShortenHandler(service ShorterService, baseURL string) *ShortenHandler {
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

	ct := req.Header.Get(HeaderContentType)

	if !strings.Contains(ct, HeaderApplicationJSONValue) && !strings.Contains(ct, HeaderApplicationXgzipValue) {
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
	if err := validation.ValidateURL(request.URL, "URL"); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID := "d1a8485a-430a-49f4-92ba-50886e1b07c6"
	shortLink, err := handler.service.CreateShortLink(req.Context(), request.URL, userID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := ShortenResponse{Result: handler.baseURL + "/" + shortLink.URL}

	res.Header().Add(HeaderContentType, HeaderApplicationJSONValue)
	if shortLink.IsDuplicated {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	err = json.NewEncoder(res).Encode(result)

	if err != nil {
		log.Zap.Errorf(ErrFailedWriteToResponse, err)
		return
	}
}
