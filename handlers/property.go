package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"myapp/middlewares"
	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/utils"
	"myapp/validate"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyHandler struct {
	Service             *services.PropertyService
	DealerService       *services.DealerService
	CloudflarePublicURL string
}

// Create
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

	dealerExists, err := h.DealerService.DealerExists(r.Context(), dealerIDObj)
	if err != nil {
		response.WithInternalError(w, r, "Failed to validate dealer: "+err.Error())
		return
	}
	if !dealerExists {
		response.WithNotFound(w, r, "Dealer not found")
		return
	}

	property.DealerID = dealerIDObj

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
		"propertyId": id.Hex(),
	})
}

// Get
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		response.WithValidationError(w, r, "Invalid property ID")
		return
	}

	property, err := h.Service.GetPropertyByID(objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WithNotFound(w, r, "Property not found")
		} else {
			response.WithInternalError(w, r, "Failed to fetch property: "+err.Error())
		}
		return
	}

	response.WithPayload(w, r, property)
}

// Update
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

	// ← PAGINATION parameters
	page := 1
	limit := 20
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	// ← FETCH properties with pagination
	properties, err := h.Service.GetPropertiesByDealer(r.Context(), dealerID, page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch properties: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WithPayload(w, r, properties)
}

func (h *PropertyHandler) GetPropertyByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	propertyNumberStr := vars["number"]

	if propertyNumberStr == "" {
		response.WithValidationError(w, r, "Property number is required")
		return
	}

	propertyNumber, err := strconv.ParseInt(propertyNumberStr, 10, 64)
	if err != nil {
		response.WithValidationError(w, r, "Invalid property number")
		return
	}

	property, err := h.Service.GetPropertyByNumber(r.Context(), propertyNumber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WithNotFound(w, r, "Property not found")
		} else {
			response.WithInternalError(w, r, "Failed to fetch property: "+err.Error())
		}
		return
	}

	response.WithPayload(w, r, property)
}

func (h *PropertyHandler) SearchProperties(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}

	// 🔹 Parse query params
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		if len(values) == 0 || values[0] == "" {
			continue
		}
		val := values[0]

		switch key {
		case "name":
			// Regex search
			filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: val, Options: "i"}}

		case "location":
			filter["location"] = val

		case "sub_location":
			filter["sub_location"] = val

		case "dealer_id":
			if objectID, err := primitive.ObjectIDFromHex(val); err == nil {
				filter["dealer_id"] = objectID
			}

		case "sold":
			// sold=true / sold=false
			if val == "true" {
				filter["sold"] = true
			} else if val == "false" {
				filter["sold"] = bson.M{"$ne": true} // either false or not set
			}

		case "min_price":
			if price, err := strconv.Atoi(val); err == nil {
				// merge with existing condition
				utils.MergeRangeCondition(filter, "sold_price", "$gte", price)
			}

		case "max_price":
			if price, err := strconv.Atoi(val); err == nil {
				utils.MergeRangeCondition(filter, "sold_price", "$lte", price)
			}

		case "status":
			// converted, available, etc.
			filter["status"] = val
		}
	}

	// 🔹 Pagination
	page := 1
	limit := 20
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	// 🔹 Fields
	fieldsParam := r.URL.Query().Get("fields")
	var fields []string
	if fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
	}

	// 🔹 Search
	properties, err := h.Service.SearchProperties(r.Context(), filter, page, limit, fields)
	if err != nil {
		http.Error(w, "Failed to search properties: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 🔹 Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
	})
}
