package services

import (
	"database/sql"
	"fmt"
	"time"

	"library-management-backend/internal/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type BookService struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewBookService(db *sql.DB, logger *logrus.Logger) *BookService {
	return &BookService{
		db:     db,
		logger: logger,
	}
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	s.logger.Info("Fetching all books")

	query := `SELECT id, title, author, year, description, isbn, genre, created_at, updated_at 
			  FROM books ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.WithError(err).Error("Failed to query books")
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year,
			&book.Description, &book.ISBN, &book.Genre,
			&book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			s.logger.WithError(err).Error("Failed to scan book")
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	s.logger.WithField("count", len(books)).Info("Successfully fetched books")
	return books, nil
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	s.logger.WithField("book_id", id).Info("Fetching book by ID")

	query := `SELECT id, title, author, year, description, isbn, genre, created_at, updated_at 
			  FROM books WHERE id = $1`

	var book models.Book
	err := s.db.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.Year,
		&book.Description, &book.ISBN, &book.Genre,
		&book.CreatedAt, &book.UpdatedAt)

	if err == sql.ErrNoRows {
		s.logger.WithField("book_id", id).Warn("Book not found")
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		s.logger.WithError(err).WithField("book_id", id).Error("Failed to fetch book")
		return nil, fmt.Errorf("failed to fetch book: %w", err)
	}

	s.logger.WithField("book_id", id).Info("Successfully fetched book")
	return &book, nil
}

func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.Book, error) {
	s.logger.WithField("title", req.Title).Info("Creating new book")

	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Author:      req.Author,
		Year:        req.Year,
		Description: req.Description,
		ISBN:        req.ISBN,
		Genre:       req.Genre,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO books (id, title, author, year, description, isbn, genre, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := s.db.Exec(query, book.ID, book.Title, book.Author, book.Year,
		book.Description, book.ISBN, book.Genre, book.CreatedAt, book.UpdatedAt)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create book")
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	s.logger.WithField("book_id", book.ID).Info("Successfully created book")
	return book, nil
}

func (s *BookService) UpdateBook(id string, req *models.UpdateBookRequest) (*models.Book, error) {
	s.logger.WithField("book_id", id).Info("Updating book")

	existingBook, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	query := `UPDATE books SET title = $1, author = $2, year = $3, description = $4, 
			  isbn = $5, genre = $6, updated_at = $7 WHERE id = $8`

	now := time.Now()
	_, err = s.db.Exec(query, req.Title, req.Author, req.Year, req.Description,
		req.ISBN, req.Genre, now, id)
	if err != nil {
		s.logger.WithError(err).WithField("book_id", id).Error("Failed to update book")
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	// Return updated book
	updatedBook := &models.Book{
		ID:          existingBook.ID,
		Title:       req.Title,
		Author:      req.Author,
		Year:        req.Year,
		Description: req.Description,
		ISBN:        req.ISBN,
		Genre:       req.Genre,
		CreatedAt:   existingBook.CreatedAt,
		UpdatedAt:   now,
	}

	s.logger.WithField("book_id", id).Info("Successfully updated book")
	return updatedBook, nil
}

func (s *BookService) DeleteBook(id string) error {
	s.logger.WithField("book_id", id).Info("Deleting book")

	result, err := s.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		s.logger.WithError(err).WithField("book_id", id).Error("Failed to delete book")
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.WithError(err).WithField("book_id", id).Error("Failed to get rows affected")
		return fmt.Errorf("failed to verify deletion: %w", err)
	}

	if rowsAffected == 0 {
		s.logger.WithField("book_id", id).Warn("Book not found for deletion")
		return fmt.Errorf("book not found")
	}

	s.logger.WithField("book_id", id).Info("Successfully deleted book")
	return nil
}
