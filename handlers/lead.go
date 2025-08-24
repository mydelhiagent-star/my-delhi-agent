package handlers

import (
	"encoding/json"
	"myapp/models"
	"myapp/services"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadHandler struct {
	Service *services.LeadService
}

func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	var lead models.Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if lead.Name == "" || lead.Phone == "" {
		http.Error(w, "Name and phone are required", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateLead(r.Context(), lead)
	if err != nil {
		http.Error(w, "Failed to create lead", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Lead created successfully",
		"lead_id": id.Hex(),
	})
}

func (h *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		http.Error(w, "Missing lead ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	lead, err := h.Service.GetLeadByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Lead not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch lead", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

func (h *LeadHandler) GetAllLeads(w http.ResponseWriter, r *http.Request) {
	leads, err := h.Service.GetAllLeads(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch leads", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leads)
}

func (h *LeadHandler) GetAllLeadsByDealerID(w http.ResponseWriter, r *http.Request) {
	dealerID := r.URL.Query().Get("dealer_id")
	if dealerID == "" {
		http.Error(w, "Missing dealer ID", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
		return
	}

	leads, err := h.Service.GetAllLeadsByDealerID(r.Context(), objID)

	if err != nil {
		http.Error(w, "Failed to fetch leads", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leads)
}

func (h *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		http.Error(w, "Missing lead ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	// Decode the fields to update
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(updateData) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateLead(r.Context(), objID, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Lead not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update lead", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Lead updated successfully",
	})
}

func (h *LeadHandler) AddPropertyInterest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["leadID"]
	if leadID == "" {
		http.Error(w, "Missing lead ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	var propertyInterest models.PropertyInterest
	if err := json.NewDecoder(r.Body).Decode(&propertyInterest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate property interest
	if propertyInterest.PropertyID.IsZero() || propertyInterest.DealerID.IsZero() {
		http.Error(w, "Property ID and dealer ID are required", http.StatusBadRequest)
		return
	}

	err = h.Service.AddPropertyInterest(r.Context(), objID, propertyInterest)
	if err != nil {
		http.Error(w, "Failed to add property interest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Property interest added successfully",
	})
}
