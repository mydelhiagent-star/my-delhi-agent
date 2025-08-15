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
	leadRouter.HandleFunc("/", h.CreateLead).Methods("POST")
	leadRouter.HandleFunc("/{id}", h.GetLead).Methods("GET")
	leadRouter.HandleFunc("/{id}", h.UpdateLead).Methods("PUT")
	// Admin-only route
	adminRouter := leadRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RequireRole("admin"))
	adminRouter.HandleFunc("/", h.GetAllLeads).Methods("GET")

}
