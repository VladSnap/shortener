package handlers

import (
	"net/http"

	"github.com/VladSnap/shortener/internal/validation"
)

// GetHandler - Обработчик запроса чтения полной ссылки по её сокращению.
type GetHandler struct {
	service ShorterService
}

// NewGetHandler - Создает новую структуру GetHandler с указателем.
func NewGetHandler(service ShorterService) *GetHandler {
	handler := new(GetHandler)
	handler.service = service
	return handler
}

// Handle - Обрабатывает входящий запрос.
func (handler *GetHandler) Handle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, ValidateErrHTTPNotGET, http.StatusBadRequest)
		return
	}

	shortID := req.PathValue("id")
	if !validation.ValidatePath(req.URL.Path) || shortID == "" {
		http.Error(res, "Request path incorrect", http.StatusBadRequest)
		return
	}

	url, err := handler.service.GetURL(req.Context(), shortID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if url == nil {
		http.Error(res, "Url not found", http.StatusNotFound)
		return
	}

	if url.IsDeleted {
		http.Error(res, "Url has been removed", http.StatusGone)
		return
	}

	res.Header().Set("Location", url.OriginalURL)
	http.Redirect(res, req, url.OriginalURL, http.StatusTemporaryRedirect)
}
