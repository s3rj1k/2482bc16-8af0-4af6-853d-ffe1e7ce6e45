package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Server represents the HTTP server with its database connection.
type Server struct {
	db  *pgxpool.Pool
	srv *http.Server
}

// New creates a new Server instance with the provided database pool and HTTP server configuration.
func New(db *pgxpool.Pool, srv *http.Server) *Server {
	return &Server{
		db:  db,
		srv: srv,
	}
}

// ListenAndServe starts the HTTP server and blocks until it's shut down.
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server and closes database connections.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// First shut down the HTTP server to stop accepting new requests
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	// Close all database connections
	s.db.Close()
	log.Println("Database connections closed")

	log.Println("Server shutdown completed")

	return nil
}
