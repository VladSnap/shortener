package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladSnap/shortener/internal/constants"
	m "github.com/VladSnap/shortener/internal/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWorker := m.NewMockDeleterWorker(ctrl)

	handler := NewDeleteHandler(mockWorker)

	t.Run("Invalid HTTP Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delete", http.NoBody)
		rec := httptest.NewRecorder()

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Http method not DELETE")
	})

	t.Run("Invalid Content-Type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/delete", http.NoBody)
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Incorrect content-type")
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		body := strings.NewReader(`invalid json`)
		req := httptest.NewRequest(http.MethodDelete, "/delete", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("Empty Short URLs", func(t *testing.T) {
		body := strings.NewReader(`[]`)
		req := httptest.NewRequest(http.MethodDelete, "/delete", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusNotAcceptable, rec.Code)
	})

	t.Run("Invalid Short URL", func(t *testing.T) {
		body := strings.NewReader(`["invalid-url"]`)
		req := httptest.NewRequest(http.MethodDelete, "/delete", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "shortURL length should be 8\n")
	})

	t.Run("Successful Request", func(t *testing.T) {
		shortURLs := []string{"gvFrtGrB", "rFvBHOug"}
		bodyBytes, _ := json.Marshal(shortURLs)
		body := bytes.NewReader(bodyBytes)

		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", body)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), constants.UserIDContextKey, "test-user-id"))

		rec := httptest.NewRecorder()

		mockWorker.EXPECT().AddToDelete(gomock.Any()).Times(1)

		handler.Handle(rec, req)

		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}
