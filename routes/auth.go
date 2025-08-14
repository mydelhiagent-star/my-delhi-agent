package routes

import (
	"myapp/handlers"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router, h *handlers.AuthHandler) {
    authRouter := r.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/signup", h.AdminSignup).Methods("POST")
    authRouter.HandleFunc("/dealers", h.CreateDealer).Methods("POST")
    authRouter.HandleFunc("/login", h.Login).Methods("POST")
}
