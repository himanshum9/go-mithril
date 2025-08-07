package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/himanshum9/go-mithril/services/auth-service/handlers"
	"github.com/himanshum9/go-mithril/services/auth-service/middleware"
)

// LoginRequest represents the login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")

	// Example of a protected endpoint (add more as needed)
	protected := r.PathPrefix("/api/protected").Subrouter()
	protected.Use(middleware.CognitoAuthMiddleware)
	protected.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaimsFromContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(claims)
	}).Methods("GET")

	log.Println("Auth service is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
