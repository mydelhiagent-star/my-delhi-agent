package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/utils"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClientHandler struct {
	Service *services.DealerClientService
}

func (h *DealerClientHandler) CreateDealerClient(w http.ResponseWriter, r *http.Request) {
	dealerID, _ := r.Context().Value(middlewares.UserIDKey).(string)
	dealerIDObj, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		response.WithValidationError(w, r, "Invalid dealer ID")
		return
	}
	var dealerClient models.DealerClient
	if err := json.NewDecoder(r.Body).Decode(&dealerClient); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	dealerClient.DealerID = dealerIDObj.Hex()

	_, err = h.Service.CreateDealerClient(r.Context(), dealerClient)
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
	})
}

func (h *DealerClientHandler) GetDealerClients(w http.ResponseWriter, r *http.Request) {
	dealerID := r.Context().Value(middlewares.UserIDKey).(string)
	_, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
		return
	}
	var params models.DealerClientQueryParams
	if err := utils.ParseQueryParams(r, &params); err != nil {
        http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
        return
    }

	params.DealerID = &dealerID

	fields := utils.ParseFieldSelection(r)

	

	dealerClients, err := h.Service.GetDealerClients(r.Context(), params, fields)
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

	var dealerClientUpdate models.DealerClientUpdate

	

	if err := json.NewDecoder(r.Body).Decode(&dealerClientUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current dealer client to fetch property_id and dealer_id
	currentClient, err := h.Service.GetDealerClientByID(r.Context(), objID.Hex())
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

	// Convert updateData to map
	updateMap := map[string]interface{}{
		"name":   updateData.Name,
		"phone":  updateData.Phone,
		"status": updateData.Status,
		"note":   updateData.Note,
	}
	err = h.Service.UpdateDealerClient(r.Context(), objID.Hex(), updateMap)
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
	err = h.Service.DeleteDealerClient(r.Context(), objID.Hex())
	if err != nil {
		http.Error(w, "Failed to delete dealer client", http.StatusInternalServerError)
		return
	}
	response.WithMessage(w, r, "Dealer client deleted successfully")
}



func (h *DealerClientHandler) CreateDealerClientPropertyInterest(w http.ResponseWriter, r *http.Request) {
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
	var dealerClientPropertyInterest models.DealerClientPropertyInterest
	if err := json.NewDecoder(r.Body).Decode(&dealerClientPropertyInterest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.Service.CreateDealerClientPropertyInterest(r.Context(), objID.Hex(), dealerClientPropertyInterest)
	if err != nil {
		http.Error(w, "Failed to create dealer client property interest", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Dealer client property interest created successfully"})
}
