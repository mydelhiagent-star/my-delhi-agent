package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/utils"

	"github.com/gorilla/mux"
)

type InquiryHandler struct {
	Service *services.InquiryService
}

func NewInquiryHandler(service *services.InquiryService) *InquiryHandler {
	return &InquiryHandler{
		Service: service,
	}
}

func (h *InquiryHandler) CreateInquiry(w http.ResponseWriter, r *http.Request) {
	var inquiry models.Inquiry
	if err := json.NewDecoder(r.Body).Decode(&inquiry); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	
	if userID := r.Context().Value(middlewares.UserIDKey); userID != nil {
		inquiry.Source = "dealer"
		if dealerID, ok := userID.(string); ok && dealerID != "" {
			inquiry.DealerID = &dealerID
		}
	} else {
		inquiry.Source = "landing_page"
	}

	createdInquiry, err := h.Service.CreateInquiry(r.Context(), inquiry)
	if err != nil {
		response.WithInternalError(w, r, "Failed to create inquiry")
		return
	}

	response.WithPayload(w, r, createdInquiry)
}

func (h *InquiryHandler) GetInquiryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	inquiry, err := h.Service.GetInquiryByID(r.Context(), id)
	if err != nil {
		response.WithNotFound(w, r, "Inquiry not found")
		return
	}

	response.WithPayload(w, r, inquiry)
}

func (h *InquiryHandler) GetAllInquiries(w http.ResponseWriter, r *http.Request) {
	var params models.InquiryQueryParams
	if err := utils.ParseQueryParams(r, &params); err != nil {
		response.WithError(w, r, "Invalid query parameters")
		return
	}

	inquiries, err := h.Service.GetAllInquiries(r.Context(), params)
	if err != nil {
		response.WithInternalError(w, r, "Failed to fetch inquiries")
		return
	}

	response.WithPayload(w, r, inquiries)
}

func (h *InquiryHandler) UpdateInquiry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates models.InquiryUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	err := h.Service.UpdateInquiry(r.Context(), id, updates)
	if err != nil {
		response.WithInternalError(w, r, "Failed to update inquiry")
		return
	}

	response.WithMessage(w, r, "Inquiry updated successfully")
}

func (h *InquiryHandler) DeleteInquiry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.Service.DeleteInquiry(r.Context(), id)
	if err != nil {
		response.WithInternalError(w, r, "Failed to delete inquiry")
		return
	}

	response.WithMessage(w, r, "Inquiry deleted successfully")
}
