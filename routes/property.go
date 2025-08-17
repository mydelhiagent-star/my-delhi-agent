package routes

import (
	
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterPropertyRoutes(r *mux.Router, h *handlers.PropertyHandler, jwtSecret string) {
    propertyRouter := r.PathPrefix("/properties").Subrouter()
	authMW := middlewares.JWTAuth(jwtSecret)
	propertyRouter.Use(authMW)
	
	
	

	dealerPropertyRouter := propertyRouter.PathPrefix("/dealer").Subrouter()
	dealerPropertyRouter.Use(middlewares.RequireRole("dealer"))
	dealerPropertyRouter.HandleFunc("/",h.CreateProperty).Methods("POST")
	dealerPropertyRouter.HandleFunc("/{id}", h.UpdateProperty).Methods("PUT")
    dealerPropertyRouter.HandleFunc("/{id}", h.DeleteProperty).Methods("DELETE")
    

	adminRouter := propertyRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RequireRole("admin"))
	adminRouter.HandleFunc("/", h.GetAllProperties).Methods("GET")

	propertyRouter.HandleFunc("/",h.GetPropertiesByDealer).Methods("GET")
}
