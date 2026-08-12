package router

import (
	"mehndi-booking-backend/internal/handler"
	"mehndi-booking-backend/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", handler.Health).Methods("GET")

	cors := middleware.NewCORS()
	r.Use(cors.Handler)

	return r
}
