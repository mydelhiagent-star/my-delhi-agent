package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/mongo_models"
	"myapp/response"
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
	if err := json.NewDecoder(r.Body).Decode(&dealerClient); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	id, err := h.Service.CreateDealerClient(r.Context(), dealerClient)
	if err != nil {
		if err.Error() == "phone number already exists" {
			response.WithConflict(w, r, "Phone number already exists")
		} else {
			response.WithInternalError(w, r, "Failed to create dealer client: "+err.Error())
		}
		return
	}

	response.WithPayload(w, r, map[string]interface{}{
		"message": "Dealer client created successfully",
		"id":      id.Hex(),
	})
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
	response.WithPayload(w, r, dealerClients)
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
	exists, err := h.Service.CheckPhoneExistsForDealer(r.Context(), currentClient.DealerID, currentClient.PropertyID, updateData.Phone)
	if err != nil {
		http.Error(w, "Failed to check phone number", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Phone number already exists", http.StatusConflict)
		return
	}

	// Convert updateData to map
	updateMap := map[string]interface{}{
		"name":   updateData.Name,
		"phone":  updateData.Phone,
		"status": updateData.Status,
		"note":   updateData.Note,
	}
	err = h.Service.UpdateDealerClient(r.Context(), objID, updateMap)
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

func (h *DealerClientHandler) UpdateDealerClientStatus(w http.ResponseWriter, r *http.Request) {
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
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.Service.UpdateDealerClientStatus(r.Context(), objID, updateData.Status)
	if err != nil {
		http.Error(w, "Failed to update dealer client status", http.StatusInternalServerError)
		return
	}
}
