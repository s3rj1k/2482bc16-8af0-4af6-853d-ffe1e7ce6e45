package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"code.local/internal/pkg/schema"
	"code.local/internal/pkg/utils"
)

// EnrollStudent handles HTTP POST requests to enroll a student in a course section.
// Validates the enrollment request, checks for conflicts and capacity via database triggers,
// creates the enrollment record, and returns the enrollment details with ID and timestamp.
func (h *Handlers) EnrollStudent(w http.ResponseWriter, r *http.Request) {
	var enrollment EnrollmentRequest

	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	query := `
		INSERT INTO enrollments (student_id, section_id)
		VALUES ($1, $2)
		RETURNING id, enrollment_date
	`

	var (
		id             int
		enrollmentDate time.Time
	)

	err := h.db.QueryRow(r.Context(), query, enrollment.StudentID, enrollment.SectionID).
		Scan(&id, &enrollmentDate)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Message {
			case "Schedule conflict detected. Cannot enroll in this section.":
				utils.SendError(w, http.StatusConflict, "Schedule conflict detected")

				return
			case "Section is full. Cannot enroll.":
				utils.SendError(w, http.StatusConflict, "Section is full")

				return
			}
		}

		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to enroll student: %v", err))

		return
	}

	result := schema.Enrollment{
		ID:             id,
		StudentID:      enrollment.StudentID,
		SectionID:      enrollment.SectionID,
		EnrollmentDate: enrollmentDate,
	}

	utils.SendJSON(w, http.StatusCreated, result)
}
