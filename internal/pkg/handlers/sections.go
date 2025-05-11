package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"

	"code.local/internal/pkg/schema"
	"code.local/internal/pkg/utils"
)

// GetSections handles HTTP GET requests to retrieve all course sections.
// Returns a list of all sections with their associated days, ordered by section ID,
// aggregating day information from the section_days table.
func (h *Handlers) GetSections(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			s.id, s.subject_id, s.teacher_id, s.classroom_id, s.section_code,
			s.start_time, s.duration_minutes, s.max_enrollment, s.current_enrollment,
			s.created_at, s.updated_at,
			ARRAY_AGG(sd.day) as days
		FROM sections s
		LEFT JOIN section_days sd ON s.id = sd.section_id
		GROUP BY s.id
		ORDER BY s.id
	`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch sections")

		return
	}
	defer rows.Close()

	var sections []schema.Section

	for rows.Next() {
		var (
			section schema.Section
			days    pq.StringArray
		)

		err := rows.Scan(
			&section.ID, &section.SubjectID, &section.TeacherID, &section.ClassroomID,
			&section.SectionCode, &section.StartTime, &section.DurationMinutes,
			&section.MaxEnrollment, &section.CurrentEnrollment,
			&section.CreatedAt, &section.UpdatedAt, &days,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan section")

			return
		}

		section.Days = []string(days)
		sections = append(sections, section)
	}

	utils.SendJSON(w, http.StatusOK, sections)
}

// CreateSection handles HTTP POST requests to create a new course section.
// Validates the section data including day values and duration constraints,
// creates the section and its associated days within a transaction,
// and returns the created section with its ID and metadata.
func (h *Handlers) CreateSection(w http.ResponseWriter, r *http.Request) {
	var sectionReq schema.CreateSectionRequest

	if err := json.NewDecoder(r.Body).Decode(&sectionReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if sectionReq.SubjectID <= 0 || sectionReq.TeacherID <= 0 || sectionReq.ClassroomID <= 0 ||
		sectionReq.SectionCode == "" || sectionReq.StartTime == "" ||
		sectionReq.DurationMinutes <= 0 || sectionReq.MaxEnrollment <= 0 ||
		len(sectionReq.Days) == 0 {
		utils.SendError(w, http.StatusBadRequest, "All fields are required")

		return
	}

	if sectionReq.DurationMinutes != 50 && sectionReq.DurationMinutes != 80 {
		utils.SendError(w, http.StatusBadRequest, "Duration minutes must be either 50 or 80")

		return
	}

	validDays := map[string]bool{
		"monday": true, "tuesday": true, "wednesday": true, "thursday": true, "friday": true,
	}

	for _, day := range sectionReq.Days {
		if !validDays[day] {
			utils.SendError(w, http.StatusBadRequest, "Days must be monday, tuesday, wednesday, thursday, or friday")

			return
		}
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to begin transaction")

		return
	}
	defer tx.Rollback(r.Context())

	var section schema.Section

	sectionQuery := `
		INSERT INTO sections (subject_id, teacher_id, classroom_id, section_code, start_time, duration_minutes, max_enrollment)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, subject_id, teacher_id, classroom_id, section_code, start_time::text, duration_minutes, max_enrollment, current_enrollment, created_at, updated_at
	`

	err = tx.QueryRow(
		r.Context(),
		sectionQuery,
		sectionReq.SubjectID,
		sectionReq.TeacherID,
		sectionReq.ClassroomID,
		sectionReq.SectionCode,
		sectionReq.StartTime,
		sectionReq.DurationMinutes,
		sectionReq.MaxEnrollment,
	).Scan(
		&section.ID, &section.SubjectID, &section.TeacherID, &section.ClassroomID,
		&section.SectionCode, &section.StartTime, &section.DurationMinutes,
		&section.MaxEnrollment, &section.CurrentEnrollment, &section.CreatedAt, &section.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Unique violation
				utils.SendError(w, http.StatusConflict, "A section with this subject and section code already exists")
			case "23514": // Check constraint violation
				utils.SendError(w, http.StatusBadRequest, "Section details violate constraints. Check time limits and duration.")
			default:
				utils.SendError(w, http.StatusInternalServerError, "Database error: "+pgErr.Message)
			}
		} else {
			utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create section: %v", err))
		}

		return
	}

	for _, day := range sectionReq.Days {
		_, err := tx.Exec(r.Context(), `
			INSERT INTO section_days (section_id, day)
			VALUES ($1, $2)
		`, section.ID, day)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to add section day: %v", err))

			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to commit transaction")

		return
	}

	section.Days = sectionReq.Days

	utils.SendJSON(w, http.StatusCreated, section)
}
