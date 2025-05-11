package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jung-kurt/gofpdf"
	"github.com/lib/pq"

	"code.local/internal/pkg/schema"
	"code.local/internal/pkg/utils"
)

// GetStudents handles HTTP GET requests to retrieve all student records.
// Returns a list of all students from the database ordered by last name, then first name.
func (h *Handlers) GetStudents(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, student_id, first_name, last_name, email, created_at, updated_at
		FROM students
		ORDER BY last_name, first_name
	`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch students")

		return
	}
	defer rows.Close()

	var students []schema.Student

	for rows.Next() {
		var student schema.Student

		err := rows.Scan(
			&student.ID, &student.StudentID, &student.FirstName,
			&student.LastName, &student.Email, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan student")

			return
		}

		students = append(students, student)
	}

	utils.SendJSON(w, http.StatusOK, students)
}

// GetStudentByID handles HTTP GET requests to retrieve a specific student by ID.
// Accepts a student ID path parameter and returns the matching student record or a not found error.
func (h *Handlers) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid student ID")

		return
	}

	var student schema.Student

	query := `
		SELECT id, student_id, first_name, last_name, email, created_at, updated_at
		FROM students
		WHERE id = $1
	`

	err = h.db.QueryRow(r.Context(), query, id).Scan(
		&student.ID, &student.StudentID, &student.FirstName,
		&student.LastName, &student.Email, &student.CreatedAt, &student.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.SendError(w, http.StatusNotFound, "Student not found")

			return
		}

		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch student")

		return
	}

	utils.SendJSON(w, http.StatusOK, student)
}

// GetStudentSchedule handles HTTP GET requests to retrieve a student's course schedule.
// Accepts a student ID path parameter and returns all courses the student is enrolled in.
func (h *Handlers) GetStudentSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid student ID")

		return
	}

	query := `
		SELECT
			section_id, subject_code, subject_name, section_code,
			teacher_first_name, teacher_last_name, building, room_number,
			start_time::text, end_time::text, duration_minutes, days
		FROM student_schedule_view
		WHERE student_id = $1
		ORDER BY subject_code, section_code
	`

	rows, err := h.db.Query(r.Context(), query, id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch schedule")

		return
	}
	defer rows.Close()

	var schedule []schema.ScheduleItem

	for rows.Next() {
		var (
			item schema.ScheduleItem
			days pq.StringArray
		)

		err := rows.Scan(
			&item.SectionID, &item.SubjectCode, &item.SubjectName, &item.SectionCode,
			&item.TeacherFirstName, &item.TeacherLastName, &item.Building, &item.RoomNumber,
			&item.StartTime, &item.EndTime, &item.DurationMinutes, &days,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan schedule item")

			return
		}

		item.Days = []string(days)
		schedule = append(schedule, item)
	}

	utils.SendJSON(w, http.StatusOK, schedule)
}

// CreateStudent handles HTTP POST requests to create a new student record.
// Validates that all required fields are provided, checks for unique constraints
// on student ID and email, and returns the created student with ID and timestamps.
func (h *Handlers) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var studentReq schema.CreateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&studentReq); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if studentReq.StudentID == "" || studentReq.FirstName == "" || studentReq.LastName == "" || studentReq.Email == "" {
		utils.SendError(w, http.StatusBadRequest, "All fields are required")

		return
	}

	var student schema.Student

	query := `
		INSERT INTO students (student_id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, student_id, first_name, last_name, email, created_at, updated_at
	`

	err := h.db.QueryRow(
		r.Context(),
		query,
		studentReq.StudentID,
		studentReq.FirstName,
		studentReq.LastName,
		studentReq.Email,
	).Scan(
		&student.ID, &student.StudentID, &student.FirstName,
		&student.LastName, &student.Email, &student.CreatedAt, &student.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			if pgErr.ConstraintName == "students_student_id_key" {
				utils.SendError(w, http.StatusConflict, "A student with this student ID already exists")
			} else if pgErr.ConstraintName == "students_email_key" {
				utils.SendError(w, http.StatusConflict, "A student with this email already exists")
			} else {
				utils.SendError(w, http.StatusConflict, "A duplicate entry exists")
			}

			return
		}

		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create student: %v", err))

		return
	}

	utils.SendJSON(w, http.StatusCreated, student)
}

// DropSection handles HTTP DELETE requests to remove a student from a section.
// Accepts student ID and section ID path parameters, removes the enrollment if it exists,
// and returns a success message or a not found error.
func (h *Handlers) DropSection(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.PathValue("student_id")
	sectionIDStr := r.PathValue("section_id")

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid student ID")

		return
	}

	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid section ID")

		return
	}

	query := `
		DELETE FROM enrollments 
		WHERE student_id = $1 AND section_id = $2
	`

	result, err := h.db.Exec(r.Context(), query, studentID, sectionID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to drop section")

		return
	}

	if result.RowsAffected() == 0 {
		utils.SendError(w, http.StatusNotFound, "Enrollment not found")

		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Section dropped successfully"})
}

// DownloadStudentSchedule handles HTTP GET requests to generate a PDF of a student's schedule.
// Accepts a student ID path parameter, retrieves student and schedule data,
// creates a formatted PDF document, and returns it as a downloadable file.
func (h *Handlers) DownloadStudentSchedule(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid student ID")

		return
	}

	// Get student info
	var student schema.Student

	err = h.db.QueryRow(r.Context(), `
		SELECT id, student_id, first_name, last_name, email
		FROM students WHERE id = $1
	`, id).Scan(&student.ID, &student.StudentID, &student.FirstName, &student.LastName, &student.Email)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch student info")

		return
	}

	// Get schedule items
	query := `
		SELECT
			section_id, subject_code, subject_name, section_code,
			teacher_first_name, teacher_last_name, building, room_number,
			start_time::text, end_time::text, duration_minutes, days
		FROM student_schedule_view
		WHERE student_id = $1
		ORDER BY days[1], start_time
	`

	rows, err := h.db.Query(r.Context(), query, id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch schedule")

		return
	}
	defer rows.Close()

	var scheduleItems []schema.ScheduleItem

	for rows.Next() {
		var (
			item schema.ScheduleItem
			days []string
		)

		err = rows.Scan(
			&item.SectionID, &item.SubjectCode, &item.SubjectName, &item.SectionCode,
			&item.TeacherFirstName, &item.TeacherLastName, &item.Building, &item.RoomNumber,
			&item.StartTime, &item.EndTime, &item.DurationMinutes, &days,
		)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to scan schedule item")

			return
		}

		item.Days = days
		scheduleItems = append(scheduleItems, item)
	}

	// Create PDF with Landscape orientation
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10) // Reduced side margins to maximize width
	pdf.AddPage()

	// Add a title with styling - black text
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(0, 0, 0) // Black title text
	title := fmt.Sprintf("Schedule for %s %s (%s)", student.FirstName, student.LastName, student.StudentID)
	pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Add a horizontal line after the title
	pdf.SetDrawColor(0, 102, 51) // Green line
	pdf.SetLineWidth(0.5)
	pageWidth := 297 - 10 - 10 // A4 width minus margins
	pdf.Line(10, pdf.GetY(), float64(10+pageWidth), pdf.GetY())
	pdf.Ln(5)

	// Calculate the full table width based on page width
	fullTableWidth := pageWidth

	// Define column widths as percentages
	colWidthPercentages := []float64{12, 20, 13, 25, 15, 15} // Total 100%
	colWidths := make([]float64, len(colWidthPercentages))
	for i, percentage := range colWidthPercentages {
		colWidths[i] = float64(fullTableWidth) * percentage / 100
	}

	headers := []string{"Days", "Time", "Subject", "Title", "Instructor", "Location"}

	// Helper function to create table cells with borders and alignment
	tableCell := func(width float64, text string, align, border string) {
		pdf.CellFormat(width, 8, text, border, 0, align, false, 0, "")
	}

	// Create styled table header
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(0, 0, 0)       // Changed to BLACK text for header instead of white
	pdf.SetFillColor(240, 240, 240) // Light gray background for header (instead of green)

	for i, header := range headers {
		tableCell(colWidths[i], header, "C", "1")
	}
	pdf.Ln(-1)

	// Set style for table content
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(0, 0, 0) // Black text for content

	for _, item := range scheduleItems {
		daysStr := strings.Join(utils.FormatDays(item.Days), ", ")
		timeStr := fmt.Sprintf("%s - %s", item.StartTime, item.EndTime)
		instructorName := fmt.Sprintf("%s %s", item.TeacherFirstName, item.TeacherLastName)
		location := fmt.Sprintf("%s %s", item.Building, item.RoomNumber)

		// Add each cell with proper formatting
		tableCell(colWidths[0], daysStr, "L", "1")
		tableCell(colWidths[1], timeStr, "L", "1")
		tableCell(colWidths[2], item.SubjectCode, "L", "1")
		tableCell(colWidths[3], item.SubjectName, "L", "1")
		tableCell(colWidths[4], instructorName, "L", "1")
		tableCell(colWidths[5], location, "L", "1")

		pdf.Ln(-1)
	}

	// Generate PDF bytes
	var buf bytes.Buffer

	err = pdf.Output(&buf)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate PDF")

		return
	}

	// Set headers for download
	fileName := fmt.Sprintf("schedule_%s_%s.pdf", student.FirstName, student.LastName)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))

	// Send PDF
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
