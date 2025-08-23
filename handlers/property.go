package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/middlewares"
	"myapp/models"
	"myapp/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyHandler struct {
	Service             *services.PropertyService
	CloudflarePublicURL string
}

// Create
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property models.Property
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		http.Error(w, "Invalid request body "+err.Error(), http.StatusBadRequest)
		return
	}

	dealerIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || dealerIDStr == "" {
		http.Error(w, "Unauthorized - Dealer ID not found", http.StatusUnauthorized)
		return
	}

	dealerIDObj, err := primitive.ObjectIDFromHex(dealerIDStr)
	if err != nil {
		http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
		return
	}

	property.DealerID = dealerIDObj

	publicURLPrefix := h.CloudflarePublicURL

	for i, photoKey := range property.Photos {
		if photoKey != "" {
			property.Photos[i] = publicURLPrefix + photoKey
		}
	}

	for i, videoKey := range property.Videos {
		if videoKey != "" {
			property.Videos[i] = publicURLPrefix + videoKey
		}
	}

	id, err := h.Service.CreateProperty(r.Context(), property)
	if err != nil {
		http.Error(w, "Failed to create property", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Property created successfully",
		"propertyId": id.Hex(),
	})
}

// Get
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	property, err := h.Service.GetPropertyByID(objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Property not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch property", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(property)
}

// Update
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	var updates models.Property
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateProperty(objID, updates); err != nil {
		http.Error(w, "Failed to update property", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Property updated successfully"})
}

// Delete
func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteProperty(objID); err != nil {
		http.Error(w, "Failed to delete property", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Property deleted successfully"})
}

func (h *PropertyHandler) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.Service.GetAllProperties(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch properties: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}

func (h *PropertyHandler) GetPropertiesByDealer(w http.ResponseWriter, r *http.Request) {
	var dealerID primitive.ObjectID
	var err error

	// Get dealer_id from query parameter
	dealerIDParam := r.URL.Query().Get("dealer_id")

	if dealerIDParam != "" {
		// Admin or external call with dealer_id in query
		dealerID, err = primitive.ObjectIDFromHex(dealerIDParam)
		if err != nil {
			http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
			return
		}
	} else {
		// No dealer_id in query → must be dealer calling their own
		role, _ := r.Context().Value(middlewares.UserRoleKey).(string)
		userIDStr, _ := r.Context().Value(middlewares.UserIDKey).(string)

		if role != "dealer" {
			http.Error(w, "Dealer ID is required", http.StatusBadRequest)
			return
		}

		dealerID, err = primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			http.Error(w, "Invalid dealer ID from token", http.StatusUnauthorized)
			return
		}
	}

	// Fetch properties
	properties, err := h.Service.GetPropertiesByDealer(r.Context(), dealerID)
	if err != nil {
		http.Error(w, "Failed to fetch properties: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}
