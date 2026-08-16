package router

import (
	"encoding/json"
	"net/http"

	"mehndi-booking-backend/internal/auth"
	"mehndi-booking-backend/internal/availability"
	"mehndi-booking-backend/internal/booking"
	"mehndi-booking-backend/internal/category"
	"mehndi-booking-backend/internal/design"
	"mehndi-booking-backend/internal/handler"
	"mehndi-booking-backend/internal/middleware"
	"mehndi-booking-backend/internal/notification"
	"mehndi-booking-backend/internal/user"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewRouter(userHandler *user.Handler, authHandler *auth.Handler, categoryHandler *category.Handler, designHandler *design.Handler, availabilityHandler *availability.Handler, bookingHandler *booking.Handler, slotHandler *booking.SlotHandler, notificationHandler *notification.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", handler.Health).Methods("GET")

	r.HandleFunc("/api/users/{id}", userHandler.GetUserProfile).Methods("GET")
	r.Handle("/api/users/me", middleware.RequireAuth(http.HandlerFunc(userHandler.UpdateMyProfile))).Methods("PATCH")

	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	r.Handle("/api/auth/logout", middleware.RequireAuth(http.HandlerFunc(authHandler.Logout))).Methods("POST")
	r.Handle("/api/auth/me", middleware.RequireAuth(http.HandlerFunc(authHandler.GetCurrentUser))).Methods("GET")
	r.Handle("/api/auth/password", middleware.RequireAuth(http.HandlerFunc(authHandler.ChangePassword))).Methods("PATCH")

	r.HandleFunc("/api/categories", categoryHandler.ListCategories).Methods("GET")
	r.HandleFunc("/api/categories/{id}", categoryHandler.GetCategory).Methods("GET")

	r.HandleFunc("/api/designs", designHandler.ListDesigns).Methods("GET")
	r.HandleFunc("/api/designs/{id}", designHandler.GetDesign).Methods("GET")
	r.HandleFunc("/api/designs/category/{category_id}", designHandler.ListDesignsByCategory).Methods("GET")
	r.HandleFunc("/api/designs/search", designHandler.SearchDesigns).Methods("GET")

	r.HandleFunc("/api/availability/slots", slotHandler.GetAvailableSlots).Methods("GET")

	r.Handle("/api/bookings", middleware.RequireAuth(http.HandlerFunc(bookingHandler.ListBookings))).Methods("GET")
	r.HandleFunc("/api/bookings", bookingHandler.CreateBooking).Methods("POST")
	r.Handle("/api/bookings/{id}", middleware.RequireAuth(http.HandlerFunc(bookingHandler.GetBooking))).Methods("GET")
	r.Handle("/api/bookings/{id}/cancel", middleware.RequireAuth(http.HandlerFunc(bookingHandler.CancelBooking))).Methods("PATCH")

	r.Handle("/api/notifications", middleware.RequireAuth(http.HandlerFunc(notificationHandler.ListNotifications))).Methods("GET")
	r.Handle("/api/notifications/unread", middleware.RequireAuth(http.HandlerFunc(notificationHandler.ListUnreadNotifications))).Methods("GET")
	r.Handle("/api/notifications/unread/count", middleware.RequireAuth(http.HandlerFunc(notificationHandler.GetUnreadCount))).Methods("GET")
	r.Handle("/api/notifications/{id}/read", middleware.RequireAuth(http.HandlerFunc(notificationHandler.MarkAsRead))).Methods("PATCH")
	r.Handle("/api/notifications/read-all", middleware.RequireAuth(http.HandlerFunc(notificationHandler.MarkAllAsRead))).Methods("PATCH")

	r.Handle("/api/admin/categories", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(categoryHandler.AdminListCategories)))).Methods("GET")
	r.HandleFunc("/api/admin/categories", categoryHandler.CreateCategory).Methods("POST")
	r.Handle("/api/admin/categories/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(categoryHandler.UpdateCategory)))).Methods("PATCH")
	r.Handle("/api/admin/categories/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(categoryHandler.DeleteCategory)))).Methods("DELETE")

	r.Handle("/api/admin/designs", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(designHandler.AdminListDesigns)))).Methods("GET")
	r.HandleFunc("/api/admin/designs", designHandler.AdminCreateDesign).Methods("POST")
	r.Handle("/api/admin/designs/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(designHandler.AdminGetDesign)))).Methods("GET")
	r.Handle("/api/admin/designs/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(designHandler.AdminUpdateDesign)))).Methods("PATCH")
	r.Handle("/api/admin/designs/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(designHandler.AdminDeleteDesign)))).Methods("DELETE")

	r.Handle("/api/admin/bookings", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(bookingHandler.AdminListBookings)))).Methods("GET")
	r.Handle("/api/admin/bookings/{id}", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(bookingHandler.AdminGetBooking)))).Methods("GET")
	r.Handle("/api/admin/bookings/{id}/status", middleware.RequireAuth(middleware.RequireRole("ADMIN")(http.HandlerFunc(bookingHandler.AdminUpdateBookingStatus)))).Methods("PATCH")

	cors := middleware.NewCORS()
	r.Use(cors.Handler)

	return r
}
