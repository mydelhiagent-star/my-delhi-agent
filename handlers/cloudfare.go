package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/middlewares"
	"myapp/response"
	"myapp/services"
	"net/http"
	"time"
)

type CloudfareHandler struct {
	Service *services.CloudflareR2Service
}

func (h *CloudfareHandler) GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	if h.Service == nil {
		http.Error(w, "Service not initialized", http.StatusInternalServerError)
		return
	}
	userID, _ := r.Context().Value(middlewares.UserIDKey).(string)

	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Count int `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 {
		http.Error(w, "Count must be greater than 0", http.StatusBadRequest)
		return
	}

	var urls []map[string]interface{}

	for i := 0; i < req.Count; i++ {
		// Generate unique filename with timestamp
		timestamp := time.Now().UnixNano()
		fileName := fmt.Sprintf("file_%d_%d", timestamp, i)

		// Create organized key structure
		key := fmt.Sprintf("users/%s/%s", userID, fileName)

		// Generate presigned URL (15 minutes expiry)
		url, err := h.Service.GeneratePresignedURL(r.Context(), key, 15*time.Minute)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to generate upload URL")
			return
		}

		urls = append(urls, map[string]interface{}{
			"presignedUrl": url, // ← CHANGED: "upload_url" → "presignedUrl"
			"fileKey":      key, // ← CHANGED: "key" → "fileKey"
		})
	}

	// ← CHANGED: Return proper response structure
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"presignedUrls": urls,
	})
}
