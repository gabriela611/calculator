package handler

import (
	"net/http"

	"backend/internal/service"
)

// CORSMiddleware wraps an http.Handler with CORS headers.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewRouter sets up routes for the application.
func NewRouter(svc service.CalculatorService) http.Handler {
	mux := http.NewServeMux()
	calcHandler := NewCalculatorHandler(svc)

	mux.HandleFunc("/api/v1/calculate", calcHandler.Calculate)

	return CORSMiddleware(mux)
}
