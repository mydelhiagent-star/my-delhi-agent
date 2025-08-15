package routes

import (
	"myapp/handlers"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router, h *handlers.DealerHandler) {
    authRouter := r.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/dealers", h.CreateDealer).Methods("POST")
    authRouter.HandleFunc("/login", h.LoginDealer).Methods("POST")
}
