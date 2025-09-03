package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterDealerClientRoutes(r *mux.Router, h *handlers.DealerClientHandler, jwtSecret string) {
	dealerClientRouter := r.PathPrefix("/dealer-clients").Subrouter()
	dealerClientRouter.Use(middlewares.JWTAuth(jwtSecret))
	dealerClientRouter.HandleFunc("/", h.CreateDealerClient).Methods("POST")
	dealerClientRouter.HandleFunc("/{propertyID}", h.GetDealerClientByPropertyID).Methods("GET")
}