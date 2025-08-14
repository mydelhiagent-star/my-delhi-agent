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

func (h *AuthHandler) AdminSignup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Role != "admin" {
		http.Error(w, "Role must be admin", http.StatusForbidden)
		return
	}

	err := h.Service.CreateUser(r.Context(),user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin created"})
}

func (h *AuthHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {
	// Token verification skipped for brevity
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.Role = "dealer"

	err := h.Service.CreateUser(r.Context(),user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	json.NewDecoder(r.Body).Decode(&creds)

	token, err := h.Service.Login(creds.Phone, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.Service.GetAllUsers(r.Context())
    if err != nil {
        http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
