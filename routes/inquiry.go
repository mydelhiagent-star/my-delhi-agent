package routes

import (
	"myapp/handlers"
	"myapp/middlewares"

	"github.com/gorilla/mux"
)

func SetupInquiryRoutes(router *mux.Router, h *handlers.InquiryHandler) {
	inquiryRouter := router.PathPrefix("/inquiries").Subrouter()

	// Apply authentication middleware
	inquiryRouter.Use(middlewares.JWTAuth("your-jwt-secret"))

	// Routes
	inquiryRouter.HandleFunc("", h.CreateInquiry).Methods("POST")
	inquiryRouter.HandleFunc("", h.GetAllInquiries).Methods("GET")
	inquiryRouter.HandleFunc("/{id}", h.GetInquiryByID).Methods("GET")
	inquiryRouter.HandleFunc("/{id}", h.UpdateInquiry).Methods("PUT")
	inquiryRouter.HandleFunc("/{id}", h.DeleteInquiry).Methods("DELETE")
}
