package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterDealerRoutes(r *mux.Router, h *handlers.DealerHandler, jwtSecret string) {

	// Public
	public := r.PathPrefix("/auth/dealers").Subrouter()
	public.HandleFunc("/register", h.CreateDealer).Methods("POST")
	public.HandleFunc("/login", h.LoginDealer).Methods("POST")

	// Dealer
	dealer := r.PathPrefix("/dealers").Subrouter()
	dealer.Use(middlewares.JWTAuth(jwtSecret))
	dealer.Use(middlewares.RequireRole("dealer"))

	// Admin
	admin := r.PathPrefix("/admin/dealers").Subrouter()
	// admin.Use(middlewares.JWTAuth(jwtSecret))
	// admin.Use(middlewares.RequireRole("admin"))
	admin.HandleFunc("/by-sublocation", h.GetDealersBySubLocation).Methods("GET")
	admin.HandleFunc("/locations/sublocations",h.GetLocationsWithSubLocations).Methods("GET")

}
