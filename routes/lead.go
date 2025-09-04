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
	leadRouter.HandleFunc("/search", h.SearchLeads).Methods("GET")
	leadRouter.HandleFunc("/{leadID}/property-details", h.GetLeadPropertyDetails).Methods("GET")

	// Admin-only route
	adminRouter := leadRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RequireRole("admin"))
	adminRouter.HandleFunc("/", h.CreateLead).Methods("POST")
	adminRouter.HandleFunc("/", h.GetAllLeads).Methods("GET")
	adminRouter.HandleFunc("/{leadID}", h.DeleteLead).Methods("DELETE")
	adminRouter.HandleFunc("/{leadID}", h.UpdateLead).Methods("PUT")
	adminRouter.HandleFunc("/{leadID}/properties", h.AddPropertyInterest).Methods("POST")
	adminRouter.HandleFunc("/conflicting-properties", h.GetConflictingProperties).Methods("GET")
	adminRouter.HandleFunc("/{leadID}/properties/{propertyID}", h.UpdatePropertyInterest).Methods("PUT")

}
