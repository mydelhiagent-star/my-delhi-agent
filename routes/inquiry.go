package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func SetupInquiryRoutes(router *mux.Router, h *handlers.InquiryHandler, jwtSecret string) {
	
	createRouter := router.PathPrefix("/inquiries").Subrouter()
	createRouter.Use(middlewares.OptionalJWTAuth(jwtSecret))
	createRouter.HandleFunc("", h.CreateInquiry).Methods("POST")

	
	inquiryRouter := router.PathPrefix("/inquiries").Subrouter()
	inquiryRouter.Use(middlewares.JWTAuth(jwtSecret))

	
	inquiryRouter.HandleFunc("", h.GetAllInquiries).Methods("GET")
	inquiryRouter.HandleFunc("/{id}", h.GetInquiryByID).Methods("GET")
	inquiryRouter.HandleFunc("/{id}", h.UpdateInquiry).Methods("PUT")
	inquiryRouter.HandleFunc("/{id}", h.DeleteInquiry).Methods("DELETE")
}
