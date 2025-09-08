package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"
	"myapp/response"
	"myapp/services"
	"myapp/validate"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerHandler struct {
	Service *services.DealerService
}

func (h *DealerHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {

	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}
	if err := validate.ValidateDealer(dealer); err != nil {
		response.WithValidationError(w, r, err.Error())
		return
	}
	err := h.Service.CreateDealer(r.Context(), dealer)
	if err != nil {
		response.WithInternalError(w, r, err.Error())
		return
	}

	response.WithMessage(w, r, "Dealer created successfully")
}

func (h *DealerHandler) LoginDealer(w http.ResponseWriter, r *http.Request) {

	var creds models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	if creds.Phone == "" || creds.Password == "" {
		response.WithValidationError(w, r, "Phone and password are required")
		return
	}

	token, err := h.Service.LoginDealer(r.Context(), creds.Phone, creds.Password)
	if err != nil {
		response.WithUnauthorized(w, r, err.Error())
		return
	}
	response.WithPayload(w, r, map[string]string{"token": token})
}

func (h *DealerHandler) GetDealersBySubLocation(w http.ResponseWriter, r *http.Request) {

	subLocation := r.URL.Query().Get("subLocation")
	if subLocation == "" {
		response.Error(w, http.StatusBadRequest, "Location is required")
		return
	}

	dealers, err := h.Service.GetDealersByLocation(r.Context(), subLocation)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch dealers: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, dealers)
}

func (h *DealerHandler) GetLocationsWithSubLocations(w http.ResponseWriter, r *http.Request) {
	result, err := h.Service.GetLocationsWithSubLocations(r.Context())

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch locations sublocations: "+err.Error())
	}

	response.WithPayload(w, r, result)
}

func (h *DealerHandler) GetDealerWithProperties(w http.ResponseWriter, r *http.Request) {
	subLocation := r.URL.Query().Get("subLocation")

	if subLocation == "" {
		response.Error(w, http.StatusBadRequest, "subLocation is required")
		return
	}

	dealerWithProps, err := h.Service.GetDealerWithProperties(r.Context(), subLocation)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch dealer with properties: "+err.Error())
		return
	}

	// If no dealer found
	if dealerWithProps == nil {
		response.Error(w, http.StatusNotFound, "No dealer found for the given subLocation")
		return
	}

	response.JSON(w, http.StatusOK, dealerWithProps)

}

func (h *DealerHandler) GetAllDealers(w http.ResponseWriter, r *http.Request) {
	dealers, err := h.Service.GetAllDealers(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch dealers: "+err.Error())
		return
	}
	response.JSON(w, http.StatusOK, dealers)
}

func (h *DealerHandler) UpdateDealer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealerID := vars["id"]

	if dealerID == "" {
		response.Error(w, http.StatusBadRequest, "Dealer ID is required")
		return
	}

	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dealerObjID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid dealer ID")
		return
	}
	if err := validate.ValidateDealer(dealer); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert dealer to map for update
	updateData := map[string]interface{}{
		"name":           dealer.Name,
		"phone":          dealer.Phone,
		"email":          dealer.Email,
		"office_address": dealer.OfficeAddress,
		"shop_name":      dealer.ShopName,
		"location":       dealer.Location,
		"sub_location":   dealer.SubLocation,
	}
	err = h.Service.UpdateDealer(r.Context(), dealerObjID, updateData)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update dealer: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Dealer updated"})

}

func (h *DealerHandler) DeleteDealer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealerID := vars["id"]

	if dealerID == "" {
		response.Error(w, http.StatusBadRequest, "Dealer ID is required")
		return
	}

	dealerObjID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid dealer ID")
		return
	}

	err = h.Service.DeleteDealer(r.Context(), dealerObjID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete dealer: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Dealer deleted"})
}

func (h *DealerHandler) ResetPasswordDealer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealerID := vars["id"]

	if dealerID == "" {
		response.Error(w, http.StatusBadRequest, "Dealer ID is required")
		return
	}

	dealerObjID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid dealer ID")
		return
	}

	// ← PARSE request body for new password
	var requestBody struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if requestBody.Password == "" {
		response.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	// ← VALIDATE password strength (optional)
	if len(requestBody.Password) < 6 {
		response.Error(w, http.StatusBadRequest, "Password must be at least 6 characters long")
		return
	}

	err = h.Service.ResetPasswordDealer(r.Context(), dealerObjID, requestBody.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to reset password: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
