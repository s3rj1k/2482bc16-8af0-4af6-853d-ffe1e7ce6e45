package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.local/internal/pkg/schema"
	"code.local/internal/pkg/utils"
)

// GetSubjects handles HTTP GET requests to retrieve all academic subjects.
// Returns a list of all subject records from the database ordered by subject code.
func (h *Handlers) GetSubjects(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, code, name, description, created_at, updated_at
		FROM subjects
		ORDER BY code
	`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch subjects")

		return
	}
	defer rows.Close()

	var subjects []schema.Subject

	for rows.Next() {
		var subject schema.Subject

		err := rows.Scan(
			&subject.ID, &subject.Code, &subject.Name,
			&subject.Description, &subject.CreatedAt, &subject.UpdatedAt,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan subject")

			return
		}

		subjects = append(subjects, subject)
	}

	utils.SendJSON(w, http.StatusOK, subjects)
}

// CreateSubject handles HTTP POST requests to create a new academic subject.
// Validates that required fields (code and name) are provided,
// creates the new subject record, and returns it with assigned ID and timestamps.
func (h *Handlers) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var subject schema.Subject

	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if subject.Code == "" || subject.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "Code and name are required")

		return
	}

	query := `
		INSERT INTO subjects (code, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := h.db.QueryRow(
		r.Context(),
		query,
		subject.Code,
		subject.Name,
		subject.Description,
	).Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create subject: %v", err))

		return
	}

	utils.SendJSON(w, http.StatusCreated, subject)
}
