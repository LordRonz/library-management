package handlers

import (
	"fmt"
	"net/http"

	"library-management-backend/internal/models"
	"library-management-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	bookService *services.BookService
	validator   *validator.Validate
	logger      *logrus.Logger
}

func NewBookHandler(bookService *services.BookService, validator *validator.Validate, logger *logrus.Logger) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		validator:   validator,
		logger:      logger,
	}
}

// @Summary Get all books
// @Description Retrieve all books from the library
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Failure 500 {object} models.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get books")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve books",
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Retrieve a specific book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Book not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to get book")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve book",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary Create a new book
// @Description Add a new book to the library
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.CreateBookRequest true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.CreateBookRequest

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

	book, err := h.bookService.CreateBook(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create book")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to create book",
		})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// @Summary Update a book
// @Description Update an existing book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body models.UpdateBookRequest true "Updated book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateBookRequest
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

	book, err := h.bookService.UpdateBook(id, &req)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Book not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to update book")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to update book",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary Delete a book
// @Description Delete a book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	err := h.bookService.DeleteBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "Book not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to delete book")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to delete book",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BookHandler) formatValidationErrors(err error) []models.ValidationError {
	var errors []models.ValidationError

	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = "This field is required"
		case "min":
			message = fmt.Sprintf("Must be at least %s characters", err.Param())
		case "max":
			message = fmt.Sprintf("Must be no more than %s characters", err.Param())
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
