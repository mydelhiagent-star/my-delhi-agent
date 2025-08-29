package handlers

import (
	"encoding/json"
	"myapp/middlewares"
	"myapp/models"
	"myapp/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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
		http.Error(w, "Failed to add property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Property added successfully",
	})
}

func (h *LeadHandler) SearchLeads(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	userRole := r.Context().Value(middlewares.UserRoleKey).(string)

	// ← BUILD filter dynamically from query parameters
	filter := bson.M{}

	if userRole == "dealer" {
		dealerID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
			return
		}
		filter["properties.dealer_id"] = dealerID
		if len(r.URL.Query()) > 0 {
			http.Error(w, "Dealers cannot use query parameters", http.StatusForbidden)
			return
		}
	}

	// Get all query parameters
	queryParams := r.URL.Query()

	for key, values := range queryParams {
		if len(values) > 0 && values[0] != "" {
			switch key {
			case "name":
				// Case-insensitive name search
				filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: values[0], Options: "i"}}
			case "phone":
				// Exact phone match
				filter["phone"] = values[0]
			case "aadhar_number":
				// Exact aadhar match
				filter["aadhar_number"] = values[0]
			case "property_id":
				// Search in properties array
				if objectID, err := primitive.ObjectIDFromHex(values[0]); err == nil {
					filter["properties.property_id"] = objectID
				}
			case "dealer_id":
				// Search in properties array
				if objectID, err := primitive.ObjectIDFromHex(values[0]); err == nil {
					filter["properties.dealer_id"] = objectID
				}
			case "status":
				// Search in properties array
				filter["properties.status"] = values[0]

			case "has_properties":
				// Check if lead has properties
				if values[0] == "true" {
					filter["properties"] = bson.M{"$exists": true, "$ne": []interface{}{}}
				} else if values[0] == "false" {
					filter["properties"] = bson.M{"$exists": false}
				}
			}
		}
	}

	// ← GET pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// ← SEARCH leads with filter and pagination
	leads, err := h.Service.SearchLeads(r.Context(), filter, page, limit)
	if err != nil {
		http.Error(w, "Failed to search leads: "+err.Error(), http.StatusInternalServerError)
		return
	}
	filteredLeads := make([]map[string]interface{}, 0)

	if userRole == "dealer" {

		for _, lead := range leads {
			// Create filtered lead with only allowed fields
			filteredLead := map[string]interface{}{
				"id":         lead.ID,
				"name":       lead.Name,
				"properties": lead.Properties,
			}

			filteredLeads = append(filteredLeads, filteredLead)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"leads": filteredLeads,
		})
		return

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"leads": leads,
	})
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

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	userRole := r.Context().Value(middlewares.UserRoleKey).(string)

	if userRole == "dealer" {
		filteredProperties := make([]models.Property, 0)
		dealerID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid dealer ID", http.StatusBadRequest)
			return
		}
		for _, property := range properties {
			if property.DealerID == dealerID {
				filteredProperties = append(filteredProperties, property)
			}
		}
		properties = filteredProperties
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}

func (h *LeadHandler) GetConflictingProperties(w http.ResponseWriter, r *http.Request) {

	conflictingProperties, err := h.Service.GetConflictingProperties(r.Context())

	if err != nil {
		http.Error(w, "Failed to fetch conflicting properties", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conflictingProperties)
}
