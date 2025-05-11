package server

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"code.local/internal/pkg/config"
)

// InitDB creates and initializes a connection pool to the PostgreSQL database.
func InitDB() (*pgxpool.Pool, error) {
	// Get database connection parameters from environment variables
	dbHost := cmp.Or(os.Getenv(config.EnvDBHost), "localhost")
	dbPort := cmp.Or(os.Getenv(config.EnvDBPort), "5432")
	dbUser := os.Getenv(config.EnvDBUser)
	dbPassword := os.Getenv(config.EnvDBPassword)
	dbName := os.Getenv(config.EnvDBName)

	// Join host and port correctly using net.JoinHostPort
	hostPort := net.JoinHostPort(dbHost, dbPort)

	// Construct the PostgreSQL connection string
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s/%s",
		dbUser,
		dbPassword,
		hostPort,
		dbName,
	)

	// Parse the connection URL into a pgxpool configuration object
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Create a new connection pool with the parsed configuration
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create a timeout context for the connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection with a ping and close the pool if it fails
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return the successfully connected pool
	return pool, nil
}
