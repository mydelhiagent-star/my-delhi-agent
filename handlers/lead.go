package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/utils"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadHandler struct {
	Service         *services.LeadService
	PropertyService *services.PropertyService
}

func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	var lead models.Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}
	if lead.Name == "" || lead.Phone == "" {
		response.WithValidationError(w, r, "Name and phone are required")
		return
	}
	id, err := h.Service.CreateLead(r.Context(), lead)
	if err != nil {
		response.WithInternalError(w, r, "Failed to create lead: "+err.Error())
		return
	}

	response.WithPayload(w, r, map[string]interface{}{
		"message": "Lead created successfully",
		"lead_id": id,
	})
}

func (h *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		response.WithValidationError(w, r, "Missing lead ID")
		return
	}

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		response.WithValidationError(w, r, "Invalid lead ID")
		return
	}

	lead, err := h.Service.GetLeadByID(r.Context(), objID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WithNotFound(w, r, "Lead not found")
		} else {
			response.WithInternalError(w, r, "Failed to fetch lead: "+err.Error())
		}
		return
	}

	response.WithPayload(w, r, lead)
}



func (h *LeadHandler) GetAllLeadsByDealerID(w http.ResponseWriter, r *http.Request) {
	dealerID := r.URL.Query().Get("dealer_id")
	if dealerID == "" {
		response.WithValidationError(w, r, "Missing dealer ID")
		return
	}
	objID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		response.WithValidationError(w, r, "Invalid dealer ID")
		return
	}

	leads, err := h.Service.GetAllLeadsByDealerID(r.Context(), objID.Hex())

	if err != nil {
		response.WithInternalError(w, r, "Failed to fetch leads: "+err.Error())
		return
	}
	response.WithPayload(w, r, leads)
}

func (h *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["leadID"]

	if leadID == "" {
		response.WithValidationError(w, r, "Missing lead ID")
		return
	}

	objID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		response.WithValidationError(w, r, "Invalid lead ID")
		return
	}

	// Decode the fields to update
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	if len(updateData) == 0 {
		response.WithValidationError(w, r, "No fields to update")
		return
	}

	err = h.Service.UpdateLead(r.Context(), objID.Hex(), updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WithNotFound(w, r, "Lead not found")
		} else {
			response.WithInternalError(w, r, "Failed to update lead: "+err.Error())
		}
		return
	}

	response.WithMessage(w, r, "Lead updated successfully")
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
	if propertyInterest.PropertyID == "" || propertyInterest.DealerID == "" {
		http.Error(w, "Property ID and dealer ID are required", http.StatusBadRequest)
		return
	}

	err = h.Service.AddPropertyInterest(r.Context(), objID.Hex(), propertyInterest)
	if err != nil {
		if err.Error() == "property already added to this lead" {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Property already added to this lead",
			})
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to add property: "+err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Property added successfully",
	})
}

func (h *LeadHandler) GetLeads(w http.ResponseWriter, r *http.Request) {
	// ← SECURITY: Safe context value extraction
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	userRole, ok := r.Context().Value(middlewares.UserRoleKey).(string)
	if !ok || userRole != "admin" {
		http.Error(w, "Unauthorized: Missing user role", http.StatusUnauthorized)
		return
	}

	var params models.LeadQueryParams
	if err := utils.ParseQueryParams(r, &params); err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}

	leads, err := h.Service.GetLeads(r.Context(), params)
	if err != nil {
		response.WithError(w, r, "Failed to fetch leads: "+err.Error())
		return
	}

	response.WithPayload(w, r, leads)	

	

	

	

	

}

func (h *LeadHandler) GetLeadPropertyDetails(w http.ResponseWriter, r *http.Request) {
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

	properties, err := h.Service.GetLeadPropertyDetails(r.Context(), objID.Hex())

	if err != nil {
		http.Error(w, "Failed to fetch lead with details", http.StatusInternalServerError)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filteredByStatus := make([]map[string]interface{}, 0)
		for _, property := range properties {
			if propStatus, ok := property["status"].(string); ok {
				if propStatus == status {
					filteredByStatus = append(filteredByStatus, property)
				}
			}
		}
		properties = filteredByStatus
	}

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	userRole := r.Context().Value(middlewares.UserRoleKey).(string)

	if userRole == "dealer" {
		filteredProperties := make([]map[string]interface{}, 0) // ← Change to []bson.M
		dealerID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
			return
		}
		for _, property := range properties {
			// ← Handle ObjectID type correctly
			if propDealerID, ok := property["dealer_id"].(primitive.ObjectID); ok {
				if propDealerID == dealerID {
					filteredProperties = append(filteredProperties, property)
				}
			}
		}
		properties = filteredProperties
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}

func (h *LeadHandler) GetPropertyDetails(w http.ResponseWriter, r *http.Request) {
	// Extract raw query params (still strings)
	soldStr := r.URL.Query().Get("sold")
	deletedStr := r.URL.Query().Get("deleted")

	// Pass to service
	propertyDetails, err := h.Service.GetPropertyDetails(r.Context(), soldStr, deletedStr)
	if err != nil {
		http.Error(w, "Failed to fetch property details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(propertyDetails)
}

func (h *LeadHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
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

	err = h.Service.DeleteLead(r.Context(), objID.Hex())

	if err != nil {
		http.Error(w, "Failed to delete lead", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Lead deleted successfully"})
}

func (h *LeadHandler) UpdatePropertyInterest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["leadID"]
	propertyID := vars["propertyID"]

	if leadID == "" || propertyID == "" {
		http.Error(w, "Missing lead ID or property ID", http.StatusBadRequest)
		return
	}

	leadObjID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}

	propertyObjID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	// Decode the status update
	var updateData struct {
		Status    string `json:"status"`
		SoldPrice int64  `json:"sold_price"`
		Note      string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := []string{"view", "ongoing", "converted", "closed"}
	isValid := false
	for _, status := range validStatuses {
		if status == updateData.Status {
			isValid = true
			break
		}
	}
	if !isValid {
		http.Error(w, "Invalid status. Must be one of: view, ongoing, converted, closed", http.StatusBadRequest)
		return
	}

	// Validate note: only allowed for "ongoing" status
	if updateData.Note != "" && updateData.Status != "ongoing" {
		http.Error(w, "Note is only allowed when status is 'ongoing'", http.StatusBadRequest)
		return
	}

	// ← CALL the specific service method
	err = h.Service.UpdatePropertyInterest(r.Context(), leadObjID.Hex(), propertyObjID.Hex(), updateData.Status, updateData.Note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Property interest not found for this lead", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update property status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if updateData.Status == "converted" {
		soldDate := time.Now()
		err = h.PropertyService.UpdateProperty(propertyObjID.Hex(), models.PropertyUpdate{
			Sold:      &[]bool{true}[0],
			SoldPrice: &updateData.SoldPrice,
			UpdatedAt: &soldDate,
			SoldDate:  &soldDate,
		})
		if err != nil {
			http.Error(w, "Failed to update property sold status", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Property status updated successfully",
	})
}


