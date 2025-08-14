package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"
	"myapp/services"
)

type AuthHandler struct {
	Service *services.AuthService
}

func (h *AuthHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) LoginDealer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if creds.Phone == "" || creds.Password == "" {
        http.Error(w, "Phone and password are required", http.StatusBadRequest)
        return
    }

	token, err := h.Service.LoginDealer(creds.Phone, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) GetAllDealers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.Service.GetAllDealers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
