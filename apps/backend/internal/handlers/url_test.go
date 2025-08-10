package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-management-backend/internal/models"
	"library-management-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupURLHandler() (*URLHandler, *gin.Engine) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	validate := validator.New()
	urlService := services.NewURLService(logger)
	urlHandler := NewURLHandler(urlService, validate, logger)

	router := gin.New()
	router.POST("/url-process", urlHandler.ProcessURL)

	return urlHandler, router
}

func TestURLHandler_ProcessURL(t *testing.T) {
	_, router := setupURLHandler()

	t.Run("success", func(t *testing.T) {
		reqBody := models.URLProcessRequest{
			URL:       "https://byfood.com/food-experiences?query=1",
			Operation: "all",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/url-process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp models.URLProcessResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "https://www.byfood.com/food-experiences", resp.ProcessedURL)
	})

	t.Run("invalid json", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/url-process", bytes.NewBuffer([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", resp.Error)
	})

	t.Run("validation error - missing url", func(t *testing.T) {
		reqBody := models.URLProcessRequest{
			Operation: "all",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/url-process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp models.ValidationErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Validation Error", resp.Error)
		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "URL", resp.Errors[0].Field)
	})

	t.Run("validation error - invalid operation", func(t *testing.T) {
		reqBody := models.URLProcessRequest{
			URL:       "https://example.com",
			Operation: "invalid",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/url-process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp models.ValidationErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Validation Error", resp.Error)
		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "Operation", resp.Errors[0].Field)
	})
}
