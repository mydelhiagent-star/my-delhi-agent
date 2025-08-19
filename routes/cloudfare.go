package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterCloudFareRoutes(r *mux.Router, h *handlers.CloudfareHandler, jwtSecret string){
	cloudfareRouter := r.PathPrefix("/cloudfare").Subrouter()
	cloudfareRouter.Use(middlewares.JWTAuth(jwtSecret))
	cloudfareRouter.HandleFunc("/generate-presigned-url", h.GeneratePresignedURL).Methods("POST")

}