package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"code.local/internal/pkg/schema"
	"code.local/internal/pkg/utils"
)

// GetTeachers handles HTTP GET requests to retrieve all teacher records.
// Returns a list of all teacher data from the database ordered by last name, then first name.
func (h *Handlers) GetTeachers(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at
		FROM teachers
		ORDER BY last_name, first_name
	`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch teachers")

		return
	}
	defer rows.Close()

	var teachers []schema.Teacher

	for rows.Next() {
		var teacher schema.Teacher

		err := rows.Scan(
			&teacher.ID, &teacher.FirstName, &teacher.LastName,
			&teacher.Email, &teacher.CreatedAt, &teacher.UpdatedAt,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan teacher")

			return
		}

		teachers = append(teachers, teacher)
	}

	utils.SendJSON(w, http.StatusOK, teachers)
}

// CreateTeacher handles HTTP POST requests to create a new teacher record.
// Validates that required fields (first name, last name, and email) are provided,
// checks for email uniqueness, and returns the created teacher with ID and timestamps.
func (h *Handlers) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher schema.Teacher

	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if teacher.FirstName == "" || teacher.LastName == "" || teacher.Email == "" {
		utils.SendError(w, http.StatusBadRequest, "First name, last name, and email are required")

		return
	}

	query := `
		INSERT INTO teachers (first_name, last_name, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := h.db.QueryRow(
		r.Context(),
		query,
		teacher.FirstName,
		teacher.LastName,
		teacher.Email,
	).Scan(&teacher.ID, &teacher.CreatedAt, &teacher.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			utils.SendError(w, http.StatusConflict, "A teacher with this email already exists")

			return
		}

		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create teacher: %v", err))

		return
	}

	utils.SendJSON(w, http.StatusCreated, teacher)
}
