package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handlers encapsulates the database connection pool for API request handlers.
type Handlers struct {
	db *pgxpool.Pool
}

// New creates a new Handlers instance with the provided database connection pool.
func New(db *pgxpool.Pool) *Handlers {
	return &Handlers{
		db: db,
	}
}

// EnrollmentRequest represents the data needed to create a new enrollment.
type EnrollmentRequest struct {
	StudentID int `json:"student_id"`
	SectionID int `json:"section_id"`
}
