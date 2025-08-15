package routes

import (
	"myapp/handlers"

	"github.com/gorilla/mux"
)

func RegisterDealerRoutes(r *mux.Router, h *handlers.DealerHandler, jwtString string) {
    authRouter := r.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/dealers", h.CreateDealer).Methods("POST")
    authRouter.HandleFunc("/login", h.LoginDealer).Methods("POST")

    
}
