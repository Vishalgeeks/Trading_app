package router

import (
	"mehndi-booking-backend/internal/handler"
	"mehndi-booking-backend/internal/middleware"
	"mehndi-booking-backend/internal/user"

	"github.com/gorilla/mux"
)

func NewRouter(userHandler *user.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", handler.Health).Methods("GET")

	r.HandleFunc("/api/users/{id}", userHandler.GetUserProfile).Methods("GET")
	r.HandleFunc("/api/users/{id}", userHandler.UpdateUserProfile).Methods("PATCH")

	cors := middleware.NewCORS()
	r.Use(cors.Handler)

	return r
}
