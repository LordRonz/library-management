package models

import (
	"time"
)

type Book struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required,min=1,max=255"`
	Author      string    `json:"author" db:"author" validate:"required,min=1,max=255"`
	Year        int       `json:"year" db:"year" validate:"required,min=1000"`
	Description *string   `json:"description,omitempty" db:"description" validate:"omitempty,max=1000"`
	ISBN        *string   `json:"isbn,omitempty" db:"isbn" validate:"omitempty,max=20"`
	Genre       *string   `json:"genre,omitempty" db:"genre" validate:"omitempty,max=100"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBookRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Author      string  `json:"author" validate:"required,min=1,max=255"`
	Year        int     `json:"year" validate:"required,min=1000"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	ISBN        *string `json:"isbn,omitempty" validate:"omitempty,max=20"`
	Genre       *string `json:"genre,omitempty" validate:"omitempty,max=100"`
}

type UpdateBookRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Author      string  `json:"author" validate:"required,min=1,max=255"`
	Year        int     `json:"year" validate:"required,min=1000"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	ISBN        *string `json:"isbn,omitempty" validate:"omitempty,max=20"`
	Genre       *string `json:"genre,omitempty" validate:"omitempty,max=100"`
}

// URL Processing models
type URLProcessRequest struct {
	URL       string `json:"url" validate:"required,url"`
	Operation string `json:"operation" validate:"required,oneof=canonical redirection all"`
}

type URLProcessResponse struct {
	ProcessedURL string `json:"processed_url"`
}