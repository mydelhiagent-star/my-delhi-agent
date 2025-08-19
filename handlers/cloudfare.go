package handlers

import (
	"encoding/json"
	"myapp/services"
	"net/http"
)

type CloudfareHandler struct {
	Service *services.CloudflareR2Service
	
}



func (h *CloudfareHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	

}