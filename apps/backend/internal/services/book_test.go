package services

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"regexp"
	"testing"
	"time"

	"library-management-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBookService_GetAllBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	service := NewBookService(db, logger)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description", "isbn", "genre", "created_at", "updated_at"}).
			AddRow("1", "The Lord of the Rings", "J.R.R. Tolkien", 1954, "Epic fantasy novel.", "978-0618640157", "Fantasy", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books ORDER BY created_at DESC")).
			WillReturnRows(rows)

		books, err := service.GetAllBooks()
		assert.NoError(t, err)
		assert.NotNil(t, books)
		assert.Len(t, books, 1)
		assert.Equal(t, "The Lord of the Rings", books[0].Title)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books ORDER BY created_at DESC")).
			WillReturnError(errors.New("db error"))

		books, err := service.GetAllBooks()
		assert.Error(t, err)
		assert.Nil(t, books)
		assert.EqualError(t, err, "failed to fetch books: db error")
	})
}

func TestBookService_GetBookByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	service := NewBookService(db, logger)
	bookID := "some-uuid"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description", "isbn", "genre", "created_at", "updated_at"}).
			AddRow(bookID, "The Hobbit", "J.R.R. Tolkien", 1937, "Fantasy novel.", "978-0618260300", "Fantasy", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnRows(rows)

		book, err := service.GetBookByID(bookID)
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, bookID, book.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnError(sql.ErrNoRows)

		book, err := service.GetBookByID(bookID)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.EqualError(t, err, "book not found")
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnError(errors.New("db error"))

		book, err := service.GetBookByID(bookID)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.EqualError(t, err, "failed to fetch book: db error")
	})
}

func TestBookService_CreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	service := NewBookService(db, logger)
	req := &models.CreateBookRequest{
		Title:  "New Book",
		Author: "Test Author",
		Year:   2024,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO books (id, title, author, year, description, isbn, genre, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")).
			WithArgs(sqlmock.AnyArg(), req.Title, req.Author, req.Year, req.Description, req.ISBN, req.Genre, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		book, err := service.CreateBook(req)
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, req.Title, book.Title)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO books (id, title, author, year, description, isbn, genre, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")).
			WithArgs(sqlmock.AnyArg(), req.Title, req.Author, req.Year, req.Description, req.ISBN, req.Genre, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		book, err := service.CreateBook(req)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.EqualError(t, err, "failed to create book: db error")
	})
}

func TestBookService_UpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	service := NewBookService(db, logger)
	bookID := "some-uuid"
	req := &models.UpdateBookRequest{
		Title:  "Updated Title",
		Author: "Updated Author",
		Year:   2025,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description", "isbn", "genre", "created_at", "updated_at"}).
			AddRow(bookID, "Old Title", "Old Author", 2024, "Old Description", "Old ISBN", "Old Genre", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET title = $1, author = $2, year = $3, description = $4, isbn = $5, genre = $6, updated_at = $7 WHERE id = $8")).
			WithArgs(req.Title, req.Author, req.Year, req.Description, req.ISBN, req.Genre, sqlmock.AnyArg(), bookID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		book, err := service.UpdateBook(bookID, req)
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, req.Title, book.Title)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnError(sql.ErrNoRows)

		book, err := service.UpdateBook(bookID, req)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.EqualError(t, err, "book not found")
	})

	t.Run("db error on update", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description", "isbn", "genre", "created_at", "updated_at"}).
			AddRow(bookID, "Old Title", "Old Author", 2024, "Old Description", "Old ISBN", "Old Genre", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, description, isbn, genre, created_at, updated_at FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET title = $1, author = $2, year = $3, description = $4, isbn = $5, genre = $6, updated_at = $7 WHERE id = $8")).
			WithArgs(req.Title, req.Author, req.Year, req.Description, req.ISBN, req.Genre, sqlmock.AnyArg(), bookID).
			WillReturnError(errors.New("db error"))

		book, err := service.UpdateBook(bookID, req)
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.EqualError(t, err, "failed to update book: db error")
	})
}

func TestBookService_DeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	service := NewBookService(db, logger)
	bookID := "some-uuid"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := service.DeleteBook(bookID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := service.DeleteBook(bookID)
		assert.Error(t, err)
		assert.EqualError(t, err, "book not found")
	})

	t.Run("db error on delete", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnError(errors.New("db error"))

		err := service.DeleteBook(bookID)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete book: db error")
	})

	t.Run("error on rows affected", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = $1")).
			WithArgs(bookID).
			WillReturnResult(driver.ResultNoRows)

		err := service.DeleteBook(bookID)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to verify deletion: no RowsAffected available after DDL statement")
	})
}
