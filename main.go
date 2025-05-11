package main

import (
	"cmp"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"code.local/internal/pkg/config"
	"code.local/internal/pkg/cors"
	"code.local/internal/pkg/handlers"
	"code.local/internal/pkg/server"
)

func main() {
	// Set up HTTP server with timeouts
	srv := &http.Server{
		Addr:         cmp.Or(os.Getenv(config.EnvBindAddrPort), ":8080"),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// Set up database connection
	pool, err := server.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Create handlers
	hObj := handlers.New(pool)

	// Create server
	srvObj := server.New(pool, srv)

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Student routes
	mux.HandleFunc("GET /api/students", hObj.GetStudents)
	mux.HandleFunc("GET /api/students/{id}", hObj.GetStudentByID)
	mux.HandleFunc("GET /api/students/{id}/schedule", hObj.GetStudentSchedule)
	mux.HandleFunc("POST /api/students", hObj.CreateStudent)
	mux.HandleFunc("GET /api/students/{id}/schedule/pdf", hObj.DownloadStudentSchedule)
	mux.HandleFunc("DELETE /api/students/{student_id}/sections/{section_id}", hObj.DropSection)

	// Teacher routes
	mux.HandleFunc("GET /api/teachers", hObj.GetTeachers)
	mux.HandleFunc("POST /api/teachers", hObj.CreateTeacher)

	// Subject routes
	mux.HandleFunc("GET /api/subjects", hObj.GetSubjects)
	mux.HandleFunc("POST /api/subjects", hObj.CreateSubject)

	// Classroom routes
	mux.HandleFunc("GET /api/classrooms", hObj.GetClassrooms)
	mux.HandleFunc("POST /api/classrooms", hObj.CreateClassroom)

	// Section routes
	mux.HandleFunc("GET /api/sections", hObj.GetSections)
	mux.HandleFunc("POST /api/sections", hObj.CreateSection)

	// Enrollment routes
	mux.HandleFunc("POST /api/enrollments", hObj.EnrollStudent)

	// Apply CORS middleware
	srv.Handler = cors.Register(mux)

	// Set up signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %q...\n", srv.Addr)

		if err := srvObj.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-done
	log.Println("Server received shutdown signal")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	// Shutdown gracefully
	if err := srvObj.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
