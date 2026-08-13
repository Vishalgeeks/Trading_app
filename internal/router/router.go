package router

import (
	"encoding/json"
	"net/http"

	"mehndi-booking-backend/internal/auth"
	"mehndi-booking-backend/internal/handler"
	"mehndi-booking-backend/internal/middleware"
	"mehndi-booking-backend/internal/user"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewRouter(userHandler *user.Handler, authHandler *auth.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", handler.Health).Methods("GET")

	r.HandleFunc("/api/users/{id}", userHandler.GetUserProfile).Methods("GET")
	r.Handle("/api/users/me", middleware.RequireAuth(http.HandlerFunc(userHandler.UpdateMyProfile))).Methods("PATCH")

	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.Handle("/api/auth/logout", middleware.RequireAuth(http.HandlerFunc(authHandler.Logout))).Methods("POST")
	r.Handle("/api/auth/me", middleware.RequireAuth(http.HandlerFunc(authHandler.GetCurrentUser))).Methods("GET")
	r.Handle("/api/auth/password", middleware.RequireAuth(http.HandlerFunc(authHandler.ChangePassword))).Methods("PATCH")

	cors := middleware.NewCORS()
	r.Use(cors.Handler)

	return r
}
