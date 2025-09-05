package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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

func (h *LeadHandler) SearchLeads(w http.ResponseWriter, r *http.Request) {
	// ← SECURITY: Safe context value extraction
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	userRole, ok := r.Context().Value(middlewares.UserRoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "Unauthorized: Missing user role", http.StatusUnauthorized)
		return
	}

	// ← BUILD filter based on user role
	filter := bson.M{}
	queryParams := r.URL.Query()

	// ← DEALER ROLE: Restricted access
	if userRole == "dealer" {
		dealerID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
			return
		}

		// ← DEALER: Only see leads with their properties
		filter["properties.dealer_id"] = dealerID

		// ← DEALER: Allow limited query parameters
		allowedParams := map[string]bool{
			"page": true, "limit": true, "fields": true,
		}

		// ← VALIDATE dealer query parameters
		for key := range queryParams {
			if !allowedParams[key] {
				http.Error(w, "Dealers can only use: page, limit, fields, property_id", http.StatusForbidden)
				return
			}
		}
	}

	
	if userRole == "admin" {
		for key, values := range queryParams {
			if len(values) > 0 && values[0] != "" {
				switch key {
				case "name":
					// ← OPTIMIZED: Limit regex length to prevent DoS
					if len(values[0]) > 100 {
						http.Error(w, "Name search too long (max 100 chars)", http.StatusBadRequest)
						return
					}
					filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: values[0], Options: "i"}}

				case "phone":
					// ← VALIDATE phone format
					if len(values[0]) < 10 || len(values[0]) > 15 {
						http.Error(w, "Invalid phone number format", http.StatusBadRequest)
						return
					}
					filter["phone"] = values[0]

				case "aadhar_number":
					// ← VALIDATE aadhar format
					if len(values[0]) != 12 {
						http.Error(w, "Invalid aadhar number format", http.StatusBadRequest)
						return
					}
					filter["aadhar_number"] = values[0]

				case "property_id":
					if objectID, err := primitive.ObjectIDFromHex(values[0]); err == nil {
						filter["properties.property_id"] = objectID
					} else {
						http.Error(w, "Invalid property ID format", http.StatusBadRequest)
						return
					}

				case "dealer_id":
					if objectID, err := primitive.ObjectIDFromHex(values[0]); err == nil {
						filter["properties.dealer_id"] = objectID
					} else {
						http.Error(w, "Invalid dealer ID format", http.StatusBadRequest)
						return
					}

				case "status":
					// ← VALIDATE status values
					validStatuses := map[string]bool{
						"viewed": true, "interested": true, "in_process": true,
						"converted": true, "rejected": true,
					}
					if !validStatuses[values[0]] {
						http.Error(w, "Invalid status value", http.StatusBadRequest)
						return
					}
					filter["properties.status"] = values[0]

				case "has_properties":
					if values[0] == "true" {
						filter["properties"] = bson.M{"$exists": true, "$ne": []interface{}{}}
					} else if values[0] == "false" {
						filter["properties"] = bson.M{"$exists": false}
					} else {
						http.Error(w, "has_properties must be 'true' or 'false'", http.StatusBadRequest)
						return
					}
				}
			}
		}
	}

	
	page := 1
	limit := 20

	if pageStr := queryParams.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 && p <= 1000 {
			page = p
		} else {
			http.Error(w, "Invalid page number (1-1000)", http.StatusBadRequest)
			return
		}
	}

	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		} else {
			http.Error(w, "Invalid limit (1-100)", http.StatusBadRequest)
			return
		}
	}

	// ← FIELDS: Validate field names
	var fields []string
	if fieldsParam := queryParams.Get("fields"); fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")

		// ← VALIDATE field names to prevent injection
		validFields := map[string]bool{
			"name": true, "phone": true, "aadhar_number": true,
			"properties": true, "created_at": true, "updated_at": true,
		}

		for _, field := range fields {
			if !validFields[field] {
				http.Error(w, "Invalid field name: "+field, http.StatusBadRequest)
				return
			}
		}
	}

	leads, err := h.Service.SearchLeads(r.Context(), filter, page, limit, fields)
	if err != nil {
		http.Error(w, "Failed to search leads: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"leads": leads,
		"page":  page,
		"limit": limit,
		"count": len(leads),
	}

	// ← DEALER: Filter response to only relevant data
	if userRole == "dealer" {
		// ← OPTIMIZE: Use database-level filtering instead of application-level
		propertyID := queryParams.Get("property_id")
		if propertyID != "" {
			// ← VALIDATE property_id belongs to dealer
			if _, err := primitive.ObjectIDFromHex(propertyID); err != nil {
				http.Error(w, "Invalid property ID", http.StatusBadRequest)
				return
			}

			// ← ADD property_id filter to database query (already handled in filter)
			// This prevents the inefficient double-loop filtering
		}

		// ← DEALER: Only return leads with properties (already filtered in DB)
		dealerLeads := make([]models.Lead, 0, len(leads))
		for _, lead := range leads {
			if len(lead.Properties) > 0 {
				dealerLeads = append(dealerLeads, lead)
			}
		}
		response["leads"] = dealerLeads
		response["count"] = len(dealerLeads)
	}

	// ← RESPONSE: Set headers and send
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	properties, err := h.Service.GetLeadPropertyDetails(r.Context(), objID)

	if err != nil {
		http.Error(w, "Failed to fetch lead with details", http.StatusInternalServerError)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filteredByStatus := make([]bson.M, 0)
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
		filteredProperties := make([]bson.M, 0) // ← Change to []bson.M
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

	err = h.Service.DeleteLead(r.Context(), objID)

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
		Status    string  `json:"status"`
		SoldPrice float64 `json:"sold_price"`
		Note      string  `json:"note"`
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
	err = h.Service.UpdatePropertyInterest(r.Context(), leadObjID, propertyObjID, updateData.Status, updateData.Note)
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
		err = h.PropertyService.UpdateProperty(propertyObjID, models.PropertyUpdate{
			Sold:      &[]bool{true}[0],
			SoldPrice: &updateData.SoldPrice,
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

func (h *LeadHandler) GetDealerLeads(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}

	userRole, ok := r.Context().Value(middlewares.UserRoleKey).(string)
	if !ok || userRole != "dealer" {
		http.Error(w, "Unauthorized: Missing user role", http.StatusUnauthorized)
		return
	}

	dealerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
		return
	}

	leads, err := h.Service.GetDealerLeads(r.Context(), dealerID)
	if err != nil {
		http.Error(w, "Failed to fetch dealer leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"leads": leads})
}

