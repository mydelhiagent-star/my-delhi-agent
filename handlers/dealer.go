package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"
	"myapp/services"
)

type DealerHandler struct {
	Service *services.DealerService
}

func (h *DealerHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if dealer.Name == "" || dealer.Phone == "" || dealer.Password == "" ||
		dealer.OfficeAddress == "" || dealer.ShopName == "" ||
		dealer.Location == "" || dealer.SubLocation == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	err := h.Service.CreateDealer(r.Context(), dealer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (h *DealerHandler) LoginDealer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    if creds.Phone == "" || creds.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Phone and password are required"})
        return
    }

	token, err := h.Service.LoginDealer(r.Context(), creds.Phone, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *DealerHandler) GetDealersByLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	dealers, err := h.Service.GetDealersByLocation(r.Context(), location)
	if err != nil {
		http.Error(w, "Failed to fetch dealers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dealers)
}



