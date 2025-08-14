package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterPropertyRoutes(r *mux.Router, h *handlers.PropertyHandler, jwtSecret string) {
	authMW := middlewares.JWTAuth(jwtSecret)
    propertyRouter := r.PathPrefix("/properties").Subrouter()
	propertyRouter.Use(authMW)

    propertyRouter.HandleFunc("/", h.CreateProperty).Methods("POST")
    propertyRouter.HandleFunc("/{id}", h.GetProperty).Methods("GET")
    propertyRouter.HandleFunc("/{id}", h.UpdateProperty).Methods("PUT")
    propertyRouter.HandleFunc("/{id}", h.DeleteProperty).Methods("DELETE")

	

	    // /properties/admin routes (admin only)
	adminRouter := propertyRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RequireRole("admin"))
	adminRouter.HandleFunc("/", h.GetAllProperties).Methods("GET")
}
