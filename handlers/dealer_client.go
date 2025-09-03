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
	dealerClient.ID = primitive.NewObjectID()
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