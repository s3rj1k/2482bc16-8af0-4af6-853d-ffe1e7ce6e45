package config

import "time"

// EnvDBHost is the environment variable name for the database host address
// Example value: "localhost" or "db.example.com"
const EnvDBHost = "DB_HOST"

// EnvDBPort is the environment variable name for the database port number
// Example value: "5432" (default PostgreSQL port)
const EnvDBPort = "DB_PORT"

// EnvDBUser is the environment variable name for the database username
// Used for authentication with the PostgreSQL database
const EnvDBUser = "DB_USER"

// EnvDBPassword is the environment variable name for the database user password
// Used for authentication with the PostgreSQL database
const EnvDBPassword = "DB_PASSWORD"

// EnvDBName is the environment variable name for the database name
// Specifies which database to connect to on the PostgreSQL server
const EnvDBName = "DB_NAME"

// EnvBindAddrPort is the environment variable name for the web server port
// Example value: ":8080" (note the colon prefix for Go's HTTP server)
const EnvBindAddrPort = "APP_PORT"

// ShutdownTimeout specifies how long to wait for server to finish processing
// requests before forcefully shutting down (30 seconds)
const ShutdownTimeout = 30 * time.Second

// ReadTimeout limits the time for reading the entire request including body
// Helps prevent slow client attacks (10 seconds)
const ReadTimeout = 10 * time.Second

// WriteTimeout limits the time for writing a response to the client
// Longer than ReadTimeout to accommodate processing time (30 seconds)
const WriteTimeout = 30 * time.Second

// IdleTimeout limits how long connections may remain idle before being closed
// Helps with connection pool management (120 seconds)
const IdleTimeout = 120 * time.Second
