package handlers

import (
	"encoding/json"
	"myapp/constants"
	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/utils"
	"myapp/validate"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyHandler struct {
	Service             *services.PropertyService
	DealerService       *services.DealerService
	CloudflarePublicURL string
}


func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property models.Property
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		response.WithError(w, r, "Invalid request body: "+err.Error())
		return
	}

	if err := validate.ValidateProperty(property); err != nil {
		response.WithValidationError(w, r, err.Error())
		return
	}

	dealerIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || dealerIDStr == "" {
		response.WithUnauthorized(w, r, "Dealer ID not found")
		return
	}

	dealerIDObj, err := primitive.ObjectIDFromHex(dealerIDStr)
	if err != nil {
		response.WithValidationError(w, r, "Invalid dealer ID")
		return
	}

	dealerExists, err := h.DealerService.DealerExists(r.Context(), dealerIDObj.Hex())
	if err != nil {
		response.WithInternalError(w, r, "Failed to validate dealer: "+err.Error())
		return
	}
	if !dealerExists {
		response.WithNotFound(w, r, "Dealer not found")
		return
	}

	property.DealerID = dealerIDObj.Hex()

	now := time.Now()
	property.CreatedAt = now

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
		response.WithInternalError(w, r, "Failed to create property: "+err.Error())
		return
	}

	response.WithPayload(w, r, map[string]string{
		"message":    "Property created successfully",
		"propertyId": id,
	})
}




func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	var updates models.PropertyUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	
    

	if err := validate.ValidatePropertyUpdate(updates); err != nil {
		response.WithValidationError(w, r, err.Error())
		return
	}

	if err := h.Service.UpdateProperty(objID.Hex(), updates); err != nil {
		http.Error(w, "Failed to update property", http.StatusInternalServerError)
		return
	}

	response.WithMessage(w, r, "Property updated successfully")
}


func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteProperty(objID.Hex()); err != nil {
		response.WithInternalError(w, r, "Failed to delete property")
		return
	}

	response.WithMessage(w, r, "Property deleted successfully")
}




func (h *PropertyHandler) GetProperties(w http.ResponseWriter, r *http.Request) {
    role, _ := r.Context().Value(middlewares.UserRoleKey).(string)
	dealerID, _ := r.Context().Value(middlewares.UserIDKey).(string)
    if !constants.IsValidRole(role) {
        http.Error(w, "Unauthorized: Missing user role", http.StatusUnauthorized)
        return
    }
    
   
    var params models.PropertyQueryParams
    if err := utils.ParseQueryParams(r, &params); err != nil {
        http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
        return
    }
	if role == constants.Dealer {
		params.DealerID = &dealerID
	}
	
    
    if params.Page == nil {
		page := 1
		params.Page = &page
	}
	if params.Limit == nil {
		limit := 20
		params.Limit = &limit
	}
    
   
    filters := utils.BuildFilters(params)
    
  
    properties, err := h.Service.GetProperties(r.Context(), filters, *params.Page, *params.Limit)
    if err != nil {
        http.Error(w, "Failed to fetch properties: "+err.Error(), http.StatusInternalServerError)
        return
    }

    response.WithPayload(w, r, properties)

	
}






