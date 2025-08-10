package services

import (
	"fmt"
	"net/url"
	"strings"

	"library-management-backend/internal/models"
	"github.com/sirupsen/logrus"
)

type URLService struct {
	logger *logrus.Logger
}

func NewURLService(logger *logrus.Logger) *URLService {
	return &URLService{
		logger: logger,
	}
}

func (s *URLService) ProcessURL(req *models.URLProcessRequest) (*models.URLProcessResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"url":       req.URL,
		"operation": req.Operation,
	}).Info("Processing URL")

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse URL")
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	var processedURL string

	switch req.Operation {
	case "canonical":
		processedURL = s.canonicalOperation(parsedURL)
	case "redirection":
		processedURL = s.redirectionOperation(parsedURL)
	case "all":
		processedURL = s.allOperations(parsedURL)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	s.logger.WithFields(logrus.Fields{
		"original_url":  req.URL,
		"processed_url": processedURL,
		"operation":     req.Operation,
	}).Info("Successfully processed URL")

	return &models.URLProcessResponse{
		ProcessedURL: processedURL,
	}, nil
}

func (s *URLService) canonicalOperation(parsedURL *url.URL) string {
	parsedURL.RawQuery = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	return parsedURL.String()
}

func (s *URLService) redirectionOperation(parsedURL *url.URL) string {
	parsedURL.Host = "www.byfood.com"
	parsedURL.Path = strings.ToLower(parsedURL.Path)
	return strings.ToLower(parsedURL.String())
}

func (s *URLService) allOperations(parsedURL *url.URL) string {
	parsedURL.Host = "www.byfood.com"
	parsedURL.Path = strings.ToLower(parsedURL.Path)
	
	parsedURL.RawQuery = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	
	return strings.ToLower(parsedURL.String())
}