package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendError sends a JSON-formatted error response with the specified HTTP status code and message.
func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// SendJSON sends a JSON-formatted response with the specified HTTP status code and data.
func SendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// FormatDays converts full day names to abbreviated forms for display purposes.
func FormatDays(days []string) []string {
	formatted := make([]string, len(days))

	dayAbbr := map[string]string{
		"monday":    "M",
		"tuesday":   "Tu",
		"wednesday": "W",
		"thursday":  "Th",
		"friday":    "F",
	}

	for i, day := range days {
		formatted[i] = dayAbbr[day]
	}

	return formatted
}
