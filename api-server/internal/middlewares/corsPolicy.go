package middlewares

import "net/http"

func SetCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// setting up CORS policy
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().
				Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				http.Error(w, "No Content", http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
