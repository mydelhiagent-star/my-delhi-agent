package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/models"
	"myapp/services"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClientHandler struct {
	Service *services.DealerClientService
}

func (h *DealerClientHandler) CreateDealerClient(w http.ResponseWriter, r *http.Request) {
	var dealerClient models.DealerClient
	err := json.NewDecoder(r.Body).Decode(&dealerClient)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if phone number already exists for this dealer
	exists, err := h.Service.CheckPhoneExistsForDealer(r.Context(), dealerClient.DealerID, dealerClient.Phone)
	if err != nil {
		http.Error(w, "Failed to check phone number", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Phone number already exists", http.StatusConflict)
		return
	}

	dealerClient.Status = "unmarked"
	_, err = h.Service.CreateDealerClient(r.Context(), dealerClient)
	if err != nil {
		http.Error(w, "Failed to create dealer client", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dealerClient)
}

func (h *DealerClientHandler) GetDealerClientByPropertyID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	propertyID := vars["propertyID"]
	if propertyID == "" {
		http.Error(w, "Missing property ID", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	dealerID := r.Context().Value(middlewares.UserIDKey).(string)
	dealerIDObj, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
		return
	}
	dealerClients, err := h.Service.GetDealerClientByPropertyID(r.Context(), dealerIDObj, objID)
	if err != nil {
		http.Error(w, "Failed to fetch dealer clients", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dealerClients)
}

func (h *DealerClientHandler) UpdateDealerClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealerClientID := vars["dealerClientID"]
	if dealerClientID == "" {
		http.Error(w, "Missing dealer client ID", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		http.Error(w, "Invalid dealer client ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current dealer client to fetch property_id and dealer_id
	currentClient, err := h.Service.GetDealerClientByID(r.Context(), objID)
	if err != nil {
		http.Error(w, "Failed to fetch dealer client", http.StatusInternalServerError)
		return
	}

	// Check if phone number already exists for this dealer (excluding current client)
	exists, err := h.Service.CheckPhoneExistsForDealer(r.Context(), currentClient.DealerID, updateData.Phone)
	if err != nil {
		http.Error(w, "Failed to check phone number", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Phone number already exists", http.StatusConflict)
		return
	}

	err = h.Service.UpdateDealerClient(r.Context(), objID, updateData)
	if err != nil {
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Dealer client updated successfully"})
}


func (h *DealerClientHandler) DeleteDealerClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealerClientID := vars["dealerClientID"]
	if dealerClientID == "" {
		http.Error(w, "Missing dealer client ID", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		http.Error(w, "Invalid dealer client ID", http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteDealerClient(r.Context(), objID)
	if err != nil {
		http.Error(w, "Failed to delete dealer client", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Dealer client deleted successfully"})
}