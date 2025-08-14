package handlers

import (
    "encoding/json"
    "net/http"

    "myapp/models"
    "myapp/services"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type PropertyHandler struct {
    Service *services.PropertyService
    
}

// Create
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
    var property models.Property
    if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    id, err := h.Service.CreateProperty(property)
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


type UploadRequest struct {
    Files []string `json:"files"` 
}



func (h *PropertyHandler) GetUploadURLsHandler(w http.ResponseWriter, r *http.Request) {
    var req UploadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	uploads, err := services.GenerateSignedUploadURLs(r.Context(),req.Files)
    if err != nil {
        http.Error(w, "Failed to generate upload URLs: "+err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(uploads)
}




