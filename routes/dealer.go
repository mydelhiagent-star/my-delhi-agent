package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterDealerRoutes(r *mux.Router, h *handlers.DealerHandler, jwtSecret string) {
   
	public := r.PathPrefix("/auth/dealers").Subrouter()
	public.HandleFunc("/register", h.CreateDealer).Methods("POST")
	public.HandleFunc("/login", h.LoginDealer).Methods("POST")

	
	private := r.PathPrefix("/dealers").Subrouter()
    authMW := middlewares.JWTAuth(jwtSecret)
	private.Use(authMW)
    private.HandleFunc("/by-location", h.GetDealersByLocation).Methods("GET")
    
	
}
