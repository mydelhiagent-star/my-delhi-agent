package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterDealerClientRoutes(r *mux.Router, h *handlers.DealerClientHandler, jwtSecret string) {
	dealerClientRouter := r.PathPrefix("/dealer-clients").Subrouter()
	dealerClientRouter.Use(middlewares.JWTAuth(jwtSecret))
	dealerClientRouter.HandleFunc("", h.CreateDealerClient).Methods("POST")
	dealerClientRouter.HandleFunc("", h.GetDealerClients).Methods("GET")
	dealerClientRouter.HandleFunc("/{dealerClientID}", h.UpdateDealerClient).Methods("PUT")
	dealerClientRouter.HandleFunc("/{dealerClientID}", h.DeleteDealerClient).Methods("DELETE")
	dealerClientRouter.HandleFunc("/{dealerClientID}/properties", h.CreateDealerClientPropertyInterest).Methods("POST")
	dealerClientRouter.HandleFunc("/{dealerClientID}/properties/{propertyInterestID}", h.UpdateDealerClientPropertyInterest).Methods("PUT")
	dealerClientRouter.HandleFunc("/{dealerClientID}/properties/{propertyInterestID}", h.DeleteDealerClientPropertyInterest).Methods("DELETE")
}