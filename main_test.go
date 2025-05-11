package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"code.local/internal/pkg/handlers"
	"code.local/internal/pkg/schema"
)

const (
	apiURL = "http://localhost:8080/api"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper functions for API calls.
func createTeacher(t *testing.T, firstName, lastName, email string) schema.Teacher {
	teacher := schema.Teacher{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	resp, err := postJSON(t, apiURL+"/teachers", teacher)
	if err != nil {
		t.Fatalf("Failed to create teacher: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdTeacher schema.Teacher

	if err := json.NewDecoder(resp.Body).Decode(&createdTeacher); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdTeacher
}

func getTeachers(t *testing.T) []schema.Teacher {
	resp, err := http.Get(apiURL + "/teachers")
	if err != nil {
		t.Fatalf("Failed to get teachers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var teachers []schema.Teacher

	if err := json.NewDecoder(resp.Body).Decode(&teachers); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return teachers
}

func createSubject(t *testing.T, code, name, description string) schema.Subject {
	subject := schema.Subject{
		Code:        code,
		Name:        name,
		Description: description,
	}

	resp, err := postJSON(t, apiURL+"/subjects", subject)
	if err != nil {
		t.Fatalf("Failed to create subject: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdSubject schema.Subject

	if err := json.NewDecoder(resp.Body).Decode(&createdSubject); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdSubject
}

func getSubjects(t *testing.T) []schema.Subject {
	resp, err := http.Get(apiURL + "/subjects")
	if err != nil {
		t.Fatalf("Failed to get subjects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var subjects []schema.Subject

	if err := json.NewDecoder(resp.Body).Decode(&subjects); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return subjects
}

func createClassroom(t *testing.T, building, roomNumber string, capacity int) schema.Classroom {
	classroom := schema.Classroom{
		Building:   building,
		RoomNumber: roomNumber,
		Capacity:   capacity,
	}

	resp, err := postJSON(t, apiURL+"/classrooms", classroom)
	if err != nil {
		t.Fatalf("Failed to create classroom: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdClassroom schema.Classroom

	if err := json.NewDecoder(resp.Body).Decode(&createdClassroom); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdClassroom
}

func getClassrooms(t *testing.T) []schema.Classroom {
	resp, err := http.Get(apiURL + "/classrooms")
	if err != nil {
		t.Fatalf("Failed to get classrooms: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var classrooms []schema.Classroom

	if err := json.NewDecoder(resp.Body).Decode(&classrooms); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return classrooms
}

func createSection(t *testing.T, req schema.CreateSectionRequest) (schema.Section, error) {
	resp, err := postJSON(t, apiURL+"/sections", req)
	if err != nil {
		t.Fatalf("Failed to create section: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)

		return schema.Section{},
			fmt.Errorf("expected status %d, got %d: %s", http.StatusCreated, resp.StatusCode, string(body))
	}

	var createdSection schema.Section

	if err := json.NewDecoder(resp.Body).Decode(&createdSection); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdSection, nil
}

func getSections(t *testing.T) []schema.Section {
	resp, err := http.Get(apiURL + "/sections")
	if err != nil {
		t.Fatalf("Failed to get sections: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var sections []schema.Section

	if err := json.NewDecoder(resp.Body).Decode(&sections); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return sections
}

func createStudent(t *testing.T, req schema.CreateStudentRequest) schema.Student {
	resp, err := postJSON(t, apiURL+"/students", req)
	if err != nil {
		t.Fatalf("Failed to create student: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdStudent schema.Student

	if err := json.NewDecoder(resp.Body).Decode(&createdStudent); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdStudent
}

func getStudents(t *testing.T) []schema.Student {
	resp, err := http.Get(apiURL + "/students")
	if err != nil {
		t.Fatalf("Failed to get students: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var students []schema.Student

	if err := json.NewDecoder(resp.Body).Decode(&students); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return students
}

func getStudentByID(t *testing.T, id int) schema.Student {
	resp, err := http.Get(fmt.Sprintf("%s/students/%d", apiURL, id))
	if err != nil {
		t.Fatalf("Failed to get student: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var student schema.Student

	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return student
}

func enrollStudent(t *testing.T, studentID, sectionID int) (schema.Enrollment, error) {
	req := handlers.EnrollmentRequest{
		StudentID: studentID,
		SectionID: sectionID,
	}

	resp, err := postJSON(t, apiURL+"/enrollments", req)
	if err != nil {
		t.Fatalf("Failed to enroll student: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return schema.Enrollment{}, fmt.Errorf("failed to decode error response: %w", err)
		}

		return schema.Enrollment{}, fmt.Errorf("error enrolling student: %s (status %d)", errResp.Error, resp.StatusCode)
	}

	var enrollment schema.Enrollment

	if err := json.NewDecoder(resp.Body).Decode(&enrollment); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return enrollment, nil
}

func getStudentSchedule(t *testing.T, studentID int) []schema.ScheduleItem {
	resp, err := http.Get(fmt.Sprintf("%s/students/%d/schedule", apiURL, studentID))
	if err != nil {
		t.Fatalf("Failed to get student schedule: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var schedule []schema.ScheduleItem

	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return schedule
}

func dropSection(t *testing.T, studentID, sectionID int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/students/%d/sections/%d", apiURL, studentID, sectionID), http.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to drop section: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func downloadSchedule(t *testing.T, studentID int) []byte {
	resp, err := http.Get(fmt.Sprintf("%s/students/%d/schedule/pdf", apiURL, studentID))
	if err != nil {
		t.Fatalf("Failed to download schedule: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read PDF data: %v", err)
	}

	return pdfData
}

// Helper to make POST requests with JSON body.
func postJSON(t *testing.T, url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	return http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func TestUniversityAPI(t *testing.T) {
	t.Log("=== University Course Scheduling API Testing ===")

	t.Run("CreateTeachers", func(t *testing.T) {
		t.Log("===== TEACHER ENDPOINTS =====")

		t.Log("Creating a new teacher...")
		teacher := createTeacher(t, "Sarah", "Johnson", "sarah.johnson@university.edu")
		t.Logf("Created teacher with ID: %d\n", teacher.ID)

		t.Log("Creating a second teacher...")
		teacher2 := createTeacher(t, "Michael", "Smith", "michael.smith@university.edu")
		t.Logf("Created teacher with ID: %d\n", teacher2.ID)

		t.Log("Getting all teachers...")
		teachers := getTeachers(t)
		t.Logf("Found %d teachers\n", len(teachers))
	})

	t.Run("CreateSubjects", func(t *testing.T) {
		t.Log("===== SUBJECT ENDPOINTS =====")

		t.Log("Creating a new subject...")
		subject := createSubject(t, "CHEM101", "General Chemistry 1", "Introduction to general chemistry principles")
		t.Logf("Created subject with ID: %d\n", subject.ID)

		t.Log("Creating a second subject...")
		subject2 := createSubject(t, "CS101", "Introduction to Computer Science", "Fundamentals of programming and computer science")
		t.Logf("Created subject with ID: %d\n", subject2.ID)

		t.Log("Getting all subjects...")
		subjects := getSubjects(t)
		t.Logf("Found %d subjects\n", len(subjects))
	})

	t.Run("CreateClassrooms", func(t *testing.T) {
		t.Log("\n===== CLASSROOM ENDPOINTS =====")

		t.Log("Creating a new classroom...")
		classroom := createClassroom(t, "Science Building", "101", 40)
		t.Logf("Created classroom with ID: %d\n", classroom.ID)

		t.Log("Creating a second classroom...")
		classroom2 := createClassroom(t, "Computer Science Building", "202", 30)
		t.Logf("Created classroom with ID: %d\n", classroom2.ID)

		t.Log("Getting all classrooms...")
		classrooms := getClassrooms(t)
		t.Logf("Found %d classrooms\n", len(classrooms))
	})

	// Get the created entities for section creation
	teachers := getTeachers(t)
	subjects := getSubjects(t)
	classrooms := getClassrooms(t)

	if len(teachers) < 2 || len(subjects) < 2 || len(classrooms) < 2 {
		t.Fatal("Missing required entities for section tests")
	}

	t.Run("CreateSections", func(t *testing.T) {
		t.Log("===== SECTION ENDPOINTS =====")

		t.Log("Creating a new section...")
		sectionReq := schema.CreateSectionRequest{
			SubjectID:       subjects[0].ID,
			TeacherID:       teachers[0].ID,
			ClassroomID:     classrooms[0].ID,
			SectionCode:     "001",
			StartTime:       "08:00:00",
			DurationMinutes: 50,
			MaxEnrollment:   30,
			Days:            []string{"monday", "wednesday", "friday"},
		}

		section, err := createSection(t, sectionReq)
		if err != nil {
			t.Fatalf("Failed to create section: %v", err)
		}

		t.Logf("Created section with ID: %d\n", section.ID)

		t.Log("Creating a second section...")
		sectionReq2 := schema.CreateSectionRequest{
			SubjectID:       subjects[1].ID,
			TeacherID:       teachers[1].ID,
			ClassroomID:     classrooms[1].ID,
			SectionCode:     "101",
			StartTime:       "10:00:00",
			DurationMinutes: 80,
			MaxEnrollment:   25,
			Days:            []string{"tuesday", "thursday"},
		}

		section2, err := createSection(t, sectionReq2)
		if err != nil {
			t.Fatalf("Failed to create second section: %v", err)
		}

		t.Logf("Created section with ID: %d\n", section2.ID)

		t.Log("Getting all sections...")
		sections := getSections(t)
		t.Logf("Found %d sections\n", len(sections))
	})

	t.Run("CreateStudents", func(t *testing.T) {
		t.Log("===== STUDENT ENDPOINTS =====")

		t.Log("Creating a new student...")
		studentReq := schema.CreateStudentRequest{
			StudentID: "2024001",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@university.edu",
		}
		student := createStudent(t, studentReq)
		t.Logf("Created student with ID: %d\n", student.ID)

		t.Log("Creating a second student...")
		studentReq2 := schema.CreateStudentRequest{
			StudentID: "2024002",
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@university.edu",
		}
		student2 := createStudent(t, studentReq2)
		t.Logf("Created student with ID: %d\n", student2.ID)

		t.Log("Getting all students...")
		students := getStudents(t)
		t.Logf("Found %d students\n", len(students))

		t.Logf("Getting student with ID %d...\n", student.ID)
		retrievedStudent := getStudentByID(t, student.ID)
		t.Logf("Retrieved student: %s %s\n", retrievedStudent.FirstName, retrievedStudent.LastName)
	})

	// Test enrollments
	students := getStudents(t)
	sections := getSections(t)

	if len(students) < 2 || len(sections) < 2 {
		t.Fatal("Missing required entities for enrollment tests")
	}

	t.Run("EnrollStudent", func(t *testing.T) {
		t.Log("===== ENROLLMENT ENDPOINTS =====")

		t.Log("Enrolling student in section...")
		enrollment, err := enrollStudent(t, students[0].ID, sections[0].ID)
		if err != nil {
			t.Fatalf("Failed to enroll student: %v", err)
		}

		t.Logf("Created enrollment with ID: %d\n", enrollment.ID)

		t.Log("Enrolling student in second section...")
		enrollment2, err := enrollStudent(t, students[0].ID, sections[1].ID)
		if err != nil {
			t.Fatalf("Failed to enroll student in second section: %v", err)
		}

		t.Logf("Created enrollment with ID: %d\n", enrollment2.ID)

		t.Logf("Getting student's schedule...\n")
		schedule := getStudentSchedule(t, students[0].ID)
		t.Logf("Student has %d courses in schedule\n", len(schedule))

		for i, item := range schedule {
			t.Logf("  %d. %s (%s) - %s %s - %s:%s - %s\n",
				i+1, item.SubjectName, item.SubjectCode,
				item.TeacherFirstName, item.TeacherLastName,
				item.Building, item.RoomNumber,
				item.StartTime)
		}
	})

	t.Run("TestConflictHandling", func(t *testing.T) {
		t.Log("===== TESTING CONFLICT HANDLING =====")

		t.Log("Creating a conflicting section...")
		conflictSectionReq := schema.CreateSectionRequest{
			SubjectID:       subjects[0].ID,
			TeacherID:       teachers[0].ID,
			ClassroomID:     classrooms[0].ID,
			SectionCode:     "002",
			StartTime:       "08:30:00", // Overlaps with the first section
			DurationMinutes: 50,
			MaxEnrollment:   30,
			Days:            []string{"monday", "wednesday", "friday"},
		}

		conflictSection, err := createSection(t, conflictSectionReq)
		if err != nil {
			t.Logf("Failed to create conflicting section as expected: %v\n", err)
		} else {
			t.Logf("Created potentially conflicting section with ID: %d\n", conflictSection.ID)

			t.Log("Attempting to enroll in conflicting section...")
			_, err = enrollStudent(t, students[0].ID, conflictSection.ID)
			if err != nil {
				t.Logf("Failed to enroll in conflicting section as expected: %v\n", err)
			} else {
				t.Errorf("Expected enrollment to fail due to scheduling conflict")
			}
		}
	})

	t.Run("TestMaxEnrollment", func(t *testing.T) {
		t.Log("===== TESTING MAX ENROLLMENT =====")

		t.Log("Creating a section with low max enrollment...")
		sectionReq := schema.CreateSectionRequest{
			SubjectID:       subjects[1].ID,
			TeacherID:       teachers[1].ID,
			ClassroomID:     classrooms[1].ID,
			SectionCode:     "201",
			StartTime:       "14:00:00",
			DurationMinutes: 50,
			MaxEnrollment:   1, // Only allow one student
			Days:            []string{"friday"},
		}

		section, err := createSection(t, sectionReq)
		if err != nil {
			t.Fatalf("Failed to create section: %v", err)
		}

		t.Logf("Created section with ID: %d\n", section.ID)

		t.Log("Enrolling first student...")
		_, err = enrollStudent(t, students[0].ID, section.ID)
		if err != nil {
			t.Fatalf("Failed to enroll first student: %v", err)
		}

		t.Log("Attempting to enroll second student (should fail)...")
		_, err = enrollStudent(t, students[1].ID, section.ID)
		if err != nil {
			t.Logf("Failed to enroll second student as expected: %v\n", err)
		} else {
			t.Errorf("Expected enrollment to fail due to max enrollment reached")
		}
	})

	t.Run("DropSection", func(t *testing.T) {
		t.Log("===== DROPPING SECTION =====")

		t.Logf("Dropping first section for student %d...\n", students[0].ID)
		dropSection(t, students[0].ID, sections[0].ID)

		t.Logf("Verifying student's updated schedule...\n")
		schedule := getStudentSchedule(t, students[0].ID)
		t.Logf("Student now has %d courses in schedule\n", len(schedule))

		for i, item := range schedule {
			t.Logf("  %d. %s (%s)\n", i+1, item.SubjectName, item.SubjectCode)
		}
	})

	t.Run("DownloadSchedule", func(t *testing.T) {
		t.Log("===== DOWNLOAD STUDENT SCHEDULE =====")

		t.Logf("Downloading schedule PDF for student %d...\n", students[0].ID)
		pdfData := downloadSchedule(t, students[0].ID)

		// Save the PDF to a file
		pdfFile := fmt.Sprintf("student_%d_schedule.pdf", students[0].ID)
		if err := os.WriteFile(pdfFile, pdfData, 0o644); err != nil {
			t.Fatalf("Failed to save PDF file: %v", err)
		}

		t.Logf("Saved PDF to %s (%d bytes)\n", pdfFile, len(pdfData))
	})

	t.Log("=== End of API Testing ===")
}

func TestDuplicateEntries(t *testing.T) {
	t.Log("===== TESTING DUPLICATE ENTRIES =====")

	// Test duplicate student ID
	t.Run("DuplicateStudentID", func(t *testing.T) {
		t.Log("Creating a student...")
		studentReq := schema.CreateStudentRequest{
			StudentID: "dup_test_001",
			FirstName: "Test",
			LastName:  "Student",
			Email:     "test.student@university.edu",
		}

		student := createStudent(t, studentReq)
		t.Logf("Created student with ID: %d\n", student.ID)

		t.Log("Attempting to create a student with the same student ID...")
		studentReq2 := schema.CreateStudentRequest{
			StudentID: "dup_test_001", // Same student ID
			FirstName: "Another",
			LastName:  "Student",
			Email:     "another.student@university.edu",
		}

		resp, err := postJSON(t, apiURL+"/students", studentReq2)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with conflict due to duplicate student ID")
		} else if resp.StatusCode == http.StatusConflict {
			t.Log("Server correctly rejected duplicate student ID")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	// Test duplicate email
	t.Run("DuplicateEmail", func(t *testing.T) {
		t.Log("Creating a student...")
		studentReq := schema.CreateStudentRequest{
			StudentID: "email_test_001",
			FirstName: "Email",
			LastName:  "Test",
			Email:     "duplicate.email@university.edu",
		}

		student := createStudent(t, studentReq)
		t.Logf("Created student with ID: %d\n", student.ID)

		t.Log("Attempting to create a student with the same email...")
		studentReq2 := schema.CreateStudentRequest{
			StudentID: "email_test_002", // Different student ID
			FirstName: "Another",
			LastName:  "Person",
			Email:     "duplicate.email@university.edu", // Same email
		}

		resp, err := postJSON(t, apiURL+"/students", studentReq2)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with conflict due to duplicate email")
		} else if resp.StatusCode == http.StatusConflict {
			t.Log("Server correctly rejected duplicate email")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	// Test duplicate classroom
	t.Run("DuplicateClassroom", func(t *testing.T) {
		t.Log("Creating a classroom...")
		classroom := createClassroom(t, "Test Building", "101", 30)
		t.Logf("Created classroom with ID: %d\n", classroom.ID)

		t.Log("Attempting to create a classroom with the same building and room...")
		classroomReq := schema.Classroom{
			Building:   "Test Building",
			RoomNumber: "101", // Same building and room
			Capacity:   40,    // Different capacity
		}

		resp, err := postJSON(t, apiURL+"/classrooms", classroomReq)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with conflict due to duplicate classroom")
		} else if resp.StatusCode == http.StatusConflict {
			t.Log("Server correctly rejected duplicate classroom")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}

func TestInvalidRequests(t *testing.T) {
	t.Log("===== TESTING INVALID REQUESTS =====")

	// Test invalid section creation - invalid days
	t.Run("InvalidSectionDays", func(t *testing.T) {
		t.Log("Creating a section with invalid days...")

		// Get required IDs
		teachers := getTeachers(t)
		subjects := getSubjects(t)
		classrooms := getClassrooms(t)

		if len(teachers) < 1 || len(subjects) < 1 || len(classrooms) < 1 {
			t.Skip("Missing required entities for invalid section test")
		}

		sectionReq := schema.CreateSectionRequest{
			SubjectID:       subjects[0].ID,
			TeacherID:       teachers[0].ID,
			ClassroomID:     classrooms[0].ID,
			SectionCode:     "invalid001",
			StartTime:       "08:00:00",
			DurationMinutes: 50,
			MaxEnrollment:   30,
			Days:            []string{"saturday"}, // Invalid day
		}

		resp, err := postJSON(t, apiURL+"/sections", sectionReq)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with bad request due to invalid days")
		} else if resp.StatusCode == http.StatusBadRequest {
			t.Log("Server correctly rejected section with invalid days")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	// Test invalid section creation - invalid duration
	t.Run("InvalidSectionDuration", func(t *testing.T) {
		t.Log("Creating a section with invalid duration...")

		// Get required IDs
		teachers := getTeachers(t)
		subjects := getSubjects(t)
		classrooms := getClassrooms(t)

		if len(teachers) < 1 || len(subjects) < 1 || len(classrooms) < 1 {
			t.Skip("Missing required entities for invalid section test")
		}

		sectionReq := schema.CreateSectionRequest{
			SubjectID:       subjects[0].ID,
			TeacherID:       teachers[0].ID,
			ClassroomID:     classrooms[0].ID,
			SectionCode:     "invalid002",
			StartTime:       "08:00:00",
			DurationMinutes: 45, // Invalid duration (not 50 or 80)
			MaxEnrollment:   30,
			Days:            []string{"monday", "wednesday", "friday"},
		}

		resp, err := postJSON(t, apiURL+"/sections", sectionReq)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with bad request due to invalid duration")
		} else if resp.StatusCode == http.StatusBadRequest {
			t.Log("Server correctly rejected section with invalid duration")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	// Test missing required fields
	t.Run("MissingRequiredFields", func(t *testing.T) {
		t.Log("Creating a student with missing required fields...")

		studentReq := schema.CreateStudentRequest{
			// Missing StudentID
			FirstName: "Missing",
			LastName:  "Fields",
			Email:     "missing.fields@university.edu",
		}

		resp, err := postJSON(t, apiURL+"/students", studentReq)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Errorf("Expected request to fail with bad request due to missing fields")
		} else if resp.StatusCode == http.StatusBadRequest {
			t.Log("Server correctly rejected student with missing fields")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}

func TestEnrollmentAndDropping(t *testing.T) {
	t.Log("===== TESTING ENROLLMENT AND DROPPING =====")

	// Test multiple enrollments and drops
	t.Run("MultipleEnrollmentsAndDrops", func(t *testing.T) {
		// Create a new student
		t.Log("Creating a new student for enrollment testing...")

		studentReq := schema.CreateStudentRequest{
			StudentID: "enroll_test_001",
			FirstName: "Enrollment",
			LastName:  "Test",
			Email:     "enrollment.test@university.edu",
		}

		student := createStudent(t, studentReq)
		t.Logf("Created student with ID: %d\n", student.ID)

		// Get all sections
		sections := getSections(t)
		if len(sections) < 3 {
			// Create more sections if needed
			teachers := getTeachers(t)
			subjects := getSubjects(t)
			classrooms := getClassrooms(t)

			if len(teachers) < 1 || len(subjects) < 1 || len(classrooms) < 1 {
				t.Skip("Missing required entities for enrollment test")
			}

			// Create a new section
			t.Log("Creating additional sections for testing...")

			for i := range 3 - len(sections) {
				sectionReq := schema.CreateSectionRequest{
					SubjectID:       subjects[0].ID,
					TeacherID:       teachers[0].ID,
					ClassroomID:     classrooms[0].ID,
					SectionCode:     fmt.Sprintf("test%03d", i+1),
					StartTime:       fmt.Sprintf("%02d:00:00", 9+i),
					DurationMinutes: 50,
					MaxEnrollment:   30,
					Days:            []string{"monday", "wednesday", "friday"},
				}

				section, err := createSection(t, sectionReq)
				if err != nil {
					t.Fatalf("Failed to create section: %v", err)
				}

				t.Logf("Created section with ID: %d\n", section.ID)
				sections = append(sections, section)
			}
		}

		// Enroll student in multiple sections
		t.Log("Enrolling student in multiple sections...")

		for i := range 3 {
			enrollment, err := enrollStudent(t, student.ID, sections[i].ID)
			if err != nil {
				t.Fatalf("Failed to enroll student in section %d: %v", sections[i].ID, err)
			}

			t.Logf("Enrolled student in section %d (enrollment ID: %d)\n", sections[i].ID, enrollment.ID)
		}

		// Check student's schedule
		t.Log("Getting student's schedule after enrollments...")

		schedule := getStudentSchedule(t, student.ID)
		t.Logf("Student has %d courses in schedule\n", len(schedule))

		if len(schedule) != 3 {
			t.Errorf("Expected 3 courses in schedule, got %d", len(schedule))
		}

		// Drop one section
		t.Log("Dropping one section...")
		dropSection(t, student.ID, sections[0].ID)

		// Check student's schedule after dropping
		t.Log("Getting student's schedule after dropping section...")

		schedule = getStudentSchedule(t, student.ID)
		t.Logf("Student has %d courses in schedule\n", len(schedule))

		if len(schedule) != 2 {
			t.Errorf("Expected 2 courses in schedule after dropping, got %d", len(schedule))
		}

		// Drop all remaining sections
		t.Log("Dropping all remaining sections...")

		for i := 1; i < 3; i++ {
			dropSection(t, student.ID, sections[i].ID)
		}

		// Check student's schedule after dropping all
		t.Log("Getting student's schedule after dropping all sections...")

		schedule = getStudentSchedule(t, student.ID)
		t.Logf("Student has %d courses in schedule\n", len(schedule))

		if len(schedule) != 0 {
			t.Errorf("Expected 0 courses in schedule after dropping all, got %d", len(schedule))
		}
	})
}

func TestNonExistentResources(t *testing.T) {
	t.Log("===== TESTING NON-EXISTENT RESOURCES =====")

	// Test getting a non-existent student
	t.Run("NonExistentStudent", func(t *testing.T) {
		t.Log("Attempting to get a non-existent student...")

		resp, err := http.Get(apiURL + "/students/9999999")
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Errorf("Expected request to fail with not found for non-existent student")
		} else if resp.StatusCode == http.StatusNotFound {
			t.Log("Server correctly reported student not found")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})

	// Test dropping a non-existent enrollment
	t.Run("NonExistentEnrollment", func(t *testing.T) {
		t.Log("Attempting to drop a non-existent enrollment...")

		// First, get an existing student
		students := getStudents(t)
		if len(students) < 1 {
			t.Skip("No students found for non-existent enrollment test")
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/students/%d/sections/9999999", apiURL, students[0].ID), http.NoBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Errorf("Expected request to fail with not found for non-existent enrollment")
		} else if resp.StatusCode == http.StatusNotFound {
			t.Log("Server correctly reported enrollment not found")
		} else {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}
