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

// GetClassrooms handles HTTP GET requests to retrieve all classrooms.
// Returns an ordered list of all classroom records from the database.
func (h *Handlers) GetClassrooms(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, building, room_number, capacity, created_at, updated_at
		FROM classrooms
		ORDER BY building, room_number
	`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch classrooms")

		return
	}
	defer rows.Close()

	var classrooms []schema.Classroom

	for rows.Next() {
		var classroom schema.Classroom

		err := rows.Scan(
			&classroom.ID, &classroom.Building, &classroom.RoomNumber,
			&classroom.Capacity, &classroom.CreatedAt, &classroom.UpdatedAt,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan classroom")

			return
		}

		classrooms = append(classrooms, classroom)
	}

	utils.SendJSON(w, http.StatusOK, classrooms)
}

// CreateClassroom handles HTTP POST requests to create a new classroom.
// Validates the input, stores the new classroom, and returns the created record with its ID.
func (h *Handlers) CreateClassroom(w http.ResponseWriter, r *http.Request) {
	var classroom schema.Classroom

	if err := json.NewDecoder(r.Body).Decode(&classroom); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if classroom.Building == "" || classroom.RoomNumber == "" || classroom.Capacity <= 0 {
		utils.SendError(w, http.StatusBadRequest, "Building, room number, and a positive capacity are required")

		return
	}

	query := `
		INSERT INTO classrooms (building, room_number, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := h.db.QueryRow(
		r.Context(),
		query,
		classroom.Building,
		classroom.RoomNumber,
		classroom.Capacity,
	).Scan(&classroom.ID, &classroom.CreatedAt, &classroom.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			utils.SendError(w, http.StatusConflict, "A classroom with this building and room number already exists")

			return
		}

		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create classroom: %v", err))

		return
	}

	utils.SendJSON(w, http.StatusCreated, classroom)
}
