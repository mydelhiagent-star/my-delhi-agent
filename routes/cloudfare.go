package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func RegisterCloudFareRoutes(r *mux.Router, h *handlers.CloudfareHandler, jwtSecret string) {
	cloudfareRouter := r.PathPrefix("/cloudfare").Subrouter()
	cloudfareRouter.Use(middlewares.JWTAuth(jwtSecret))
	cloudfareRouter.HandleFunc("/presigned-urls", h.GeneratePresignedURL).Methods("POST")

}
