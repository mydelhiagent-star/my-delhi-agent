package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

// routes/property.go - HIERARCHICAL ACCESS
func RegisterPropertyRoutes(r *mux.Router, h *handlers.PropertyHandler, jwtSecret string) {
    propertyRouter := r.PathPrefix("/properties").Subrouter()
    propertyRouter.Use(middlewares.JWTAuth(jwtSecret))
    
    // ✅ Single endpoint with role-based access control
    propertyRouter.HandleFunc("", h.GetProperties).Methods("GET")
    
    // ✅ Role-specific operations
    dealerRouter := propertyRouter.PathPrefix("/dealer").Subrouter()
    dealerRouter.Use(middlewares.RequireRole("dealer"))
    dealerRouter.HandleFunc("", h.CreateProperty).Methods("POST")
    dealerRouter.HandleFunc("/{id}", h.UpdateProperty).Methods("PUT")
    dealerRouter.HandleFunc("/{id}", h.DeleteProperty).Methods("DELETE")
}
