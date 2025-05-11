package cors

import (
	"net/http"
	"strings"
)

// Register adds CORS middleware to the provided handler to enable cross-origin requests.
func Register(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Define allowed HTTP methods for cross-origin requests
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
				http.MethodOptions,
			},
			", ",
		))

		// Define allowed headers for cross-origin requests
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests immediately with a 200 OK response
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		// Process the actual request with the wrapped handler
		next.ServeHTTP(w, r)
	})
}
