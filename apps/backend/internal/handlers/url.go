package handlers

import (
	"net/http"

	"library-management-backend/internal/models"
	"library-management-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type URLHandler struct {
	urlService *services.URLService
	validator  *validator.Validate
	logger     *logrus.Logger
}

func NewURLHandler(urlService *services.URLService, validator *validator.Validate, logger *logrus.Logger) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		validator:  validator,
		logger:     logger,
	}
}

// @Summary Process URL
// @Description Process URL based on operation type (canonical, redirection, or all)
// @Tags url
// @Accept json
// @Produce json
// @Param request body models.URLProcessRequest true "URL processing request"
// @Success 200 {object} models.URLProcessResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /url-process [post]
func (h *URLHandler) ProcessURL(c *gin.Context) {
	var req models.URLProcessRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid JSON format",
		})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Error:  "Validation Error",
			Errors: validationErrors,
		})
		return
	}

	result, err := h.urlService.ProcessURL(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to process URL")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to process URL",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *URLHandler) formatValidationErrors(err error) []models.ValidationError {
	var errors []models.ValidationError
	
	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = "This field is required"
		case "url":
			message = "Must be a valid URL"
		case "oneof":
			message = "Must be one of: canonical, redirection, all"
		default:
			message = "Invalid value"
		}
		
		errors = append(errors, models.ValidationError{
			Field:   err.Field(),
			Message: message,
		})
	}
	
	return errors
}