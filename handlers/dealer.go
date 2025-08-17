package handlers

import (
	"encoding/json"
	"net/http"

	"myapp/models"
	"myapp/response"
	"myapp/services"
)

type DealerHandler struct {
	Service *services.DealerService
}

func (h *DealerHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {
	

	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if dealer.Name == "" || dealer.Phone == "" || dealer.Password == "" ||
		dealer.OfficeAddress == "" || dealer.ShopName == "" ||
		dealer.Location == "" || dealer.SubLocation == "" {
		response.Error(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	err := h.Service.CreateDealer(r.Context(), dealer)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.JSON(w, http.StatusCreated, map[string]string{"message": "Dealer created"})
}

func (h *DealerHandler) LoginDealer(w http.ResponseWriter, r *http.Request) {
	

	var creds models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if creds.Phone == "" || creds.Password == "" {
        response.Error(w, http.StatusBadRequest, "Phone and password are required")
        return
    }

	token, err := h.Service.LoginDealer(r.Context(), creds.Phone, creds.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"token": token})
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

func (h *DealerHandler) GetLocationsWithSubLocations(w http.ResponseWriter, r *http.Request){
	result,err := h.Service.GetLocationsWithSubLocations(r.Context())
	
	if err != nil {
		response.Error(w,http.StatusInternalServerError,"Failed to fetch locations sublocations: "+err.Error())
	}

	response.JSON(w,http.StatusOK,result)
}

func (h *DealerHandler) GetDealerWithProperties(w http.ResponseWriter, r *http.Request){
    subLocation := r.URL.Query().Get("subLocation")

	if subLocation == ""{
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


