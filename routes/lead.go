package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterLeadRoutes(r *mux.Router, h *handlers.LeadHandler, jwtSecret string) {
	authMW := middlewares.JWTAuth(jwtSecret)
	leadRouter := r.PathPrefix("/leads").Subrouter()
	leadRouter.Use(authMW)

	

	// Admin-only route
	adminRouter := leadRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RequireRole("admin"))
	adminRouter.HandleFunc("/", h.CreateLead).Methods("POST")
	adminRouter.HandleFunc("/", h.GetAllLeads).Methods("GET")
	adminRouter.HandleFunc("/{leadID}/properties", h.AddPropertyInterest).Methods("POST")
	adminRouter.HandleFunc("/search", h.SearchLeads).Methods("GET")
	adminRouter.HandleFunc("/{leadID}/property-details", h.GetLeadPropertyDetails).Methods("GET")


}
