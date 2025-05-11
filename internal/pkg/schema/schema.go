package schema

import "time"

// Student represents a university student record with identification and contact information.
type Student struct {
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	StudentID string    `json:"student_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	ID        int       `json:"id"`
}

// Section represents a course section with scheduling and capacity information.
type Section struct {
	CreatedAt         time.Time `json:"created_at,omitzero"`
	UpdatedAt         time.Time `json:"updated_at,omitzero"`
	SectionCode       string    `json:"section_code"`
	StartTime         string    `json:"start_time"`
	Days              []string  `json:"days"`
	ID                int       `json:"id"`
	SubjectID         int       `json:"subject_id"`
	TeacherID         int       `json:"teacher_id"`
	ClassroomID       int       `json:"classroom_id"`
	DurationMinutes   int       `json:"duration_minutes"`
	MaxEnrollment     int       `json:"max_enrollment"`
	CurrentEnrollment int       `json:"current_enrollment"`
}

// Enrollment represents a student's registration in a specific course section.
type Enrollment struct {
	EnrollmentDate time.Time `json:"enrollment_date,omitzero"`
	CreatedAt      time.Time `json:"created_at,omitzero"`
	ID             int       `json:"id"`
	StudentID      int       `json:"student_id"`
	SectionID      int       `json:"section_id"`
}

// ScheduleItem represents a course in a student's schedule with all relevant details.
type ScheduleItem struct {
	SubjectCode      string   `json:"subject_code"`
	SubjectName      string   `json:"subject_name"`
	SectionCode      string   `json:"section_code"`
	TeacherFirstName string   `json:"teacher_first_name"`
	TeacherLastName  string   `json:"teacher_last_name"`
	Building         string   `json:"building"`
	RoomNumber       string   `json:"room_number"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	Days             []string `json:"days"`
	SectionID        int      `json:"section_id"`
	DurationMinutes  int      `json:"duration_minutes"`
}

// Teacher represents a faculty member with identification and contact information.
type Teacher struct {
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	ID        int       `json:"id"`
}

// Subject represents an academic course with its code and description.
type Subject struct {
	CreatedAt   time.Time `json:"created_at,omitzero"`
	UpdatedAt   time.Time `json:"updated_at,omitzero"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ID          int       `json:"id"`
}

// Classroom represents a physical location where classes are held.
type Classroom struct {
	CreatedAt  time.Time `json:"created_at,omitzero"`
	UpdatedAt  time.Time `json:"updated_at,omitzero"`
	Building   string    `json:"building"`
	RoomNumber string    `json:"room_number"`
	ID         int       `json:"id"`
	Capacity   int       `json:"capacity"`
}

// CreateSectionRequest contains all data needed to create a new course section.
type CreateSectionRequest struct {
	SectionCode     string   `json:"section_code"`
	StartTime       string   `json:"start_time"`
	Days            []string `json:"days"`
	SubjectID       int      `json:"subject_id"`
	TeacherID       int      `json:"teacher_id"`
	ClassroomID     int      `json:"classroom_id"`
	DurationMinutes int      `json:"duration_minutes"`
	MaxEnrollment   int      `json:"max_enrollment"`
}

// CreateStudentRequest contains all data needed to create a new student record.
type CreateStudentRequest struct {
	StudentID string `json:"student_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
