package services

import (
	"io"
	"net/url"
	"testing"

	"library-management-backend/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestURLService_ProcessURL(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Discard logs for testing

	service := NewURLService(logger)

	testCases := []struct {
		name          string
		input         *models.URLProcessRequest
		expectedURL   string
		expectError   bool
		expectedError string
	}{
		{
			name: "Canonical Operation - removes query and trailing slash",
			input: &models.URLProcessRequest{
				URL:       "https://example.com/path/?query=123",
				Operation: "canonical",
			},
			expectedURL: "https://example.com/path",
			expectError: false,
		},
		{
			name: "Canonical Operation - no query or trailing slash",
			input: &models.URLProcessRequest{
				URL:       "https://example.com/path",
				Operation: "canonical",
			},
			expectedURL: "https://example.com/path",
			expectError: false,
		},
		{
			name: "Redirection Operation",
			input: &models.URLProcessRequest{
				URL:       "https://ANY.com/SOME/PATH?Q=V",
				Operation: "redirection",
			},
			expectedURL: "https://www.byfood.com/some/path?q=v",
			expectError: false,
		},
		{
			name: "Redirection Operation - from example",
			input: &models.URLProcessRequest{
				URL:       "https://BYFOOD.com/food-EXPeriences?query=abc/",
				Operation: "redirection",
			},
			expectedURL: "https://www.byfood.com/food-experiences?query=abc/",
			expectError: false,
		},
		{
			name: "All Operations - from example",
			input: &models.URLProcessRequest{
				URL:       "https://BYFOOD.com/food-EXPeriences?query=abc/",
				Operation: "all",
			},
			expectedURL: "https://www.byfood.com/food-experiences",
			expectError: false,
		},
		{
			name: "All Operations - no query params",
			input: &models.URLProcessRequest{
				URL:       "https://byfood.com/FOOD-EXPeriences/",
				Operation: "all",
			},
			expectedURL: "https://www.byfood.com/food-experiences",
			expectError: false,
		},
		{
			name: "Invalid URL",
			input: &models.URLProcessRequest{
				URL:       "://invalid-url",
				Operation: "all",
			},
			expectError:   true,
			expectedError: "invalid URL format: parse \"://invalid-url\": missing protocol scheme",
		},
		{
			name: "Unsupported Operation",
			input: &models.URLProcessRequest{
				URL:       "https://example.com",
				Operation: "unsupported",
			},
			expectError:   true,
			expectedError: "unsupported operation: unsupported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := service.ProcessURL(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tc.expectedURL, response.ProcessedURL)
			}
		})
	}
}

func TestURLService_canonicalOperation(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	service := NewURLService(logger)
	parsedURL, _ := url.Parse("https://example.com/path/subpath?param1=val1&param2=val2#fragment")
	expected := "https://example.com/path/subpath"
	result := service.canonicalOperation(parsedURL)
	assert.Equal(t, expected, result)

	parsedURL, _ = url.Parse("https://example.com/path/subpath/?param1=val1")
	expected = "https://example.com/path/subpath"
	result = service.canonicalOperation(parsedURL)
	assert.Equal(t, expected, result)
}

func TestURLService_redirectionOperation(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	service := NewURLService(logger)
	parsedURL, _ := url.Parse("https://EXAMPLE.com/PATH/SUBPATH?param1=VAL1")
	expected := "https://www.byfood.com/path/subpath?param1=val1"
	result := service.redirectionOperation(parsedURL)
	assert.Equal(t, expected, result)
}

func TestURLService_allOperations(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	service := NewURLService(logger)
	parsedURL, _ := url.Parse("https://EXAMPLE.com/PATH/SUBPATH/?param1=VAL1")
	expected := "https://www.byfood.com/path/subpath"
	result := service.allOperations(parsedURL)
	assert.Equal(t, expected, result)
}
