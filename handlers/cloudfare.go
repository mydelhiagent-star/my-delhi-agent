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
		response.WithInternalError(w, r, "Service not initialized")
		return
	}
	userID, _ := r.Context().Value(middlewares.UserIDKey).(string)

	if userID == "" {
		response.WithUnauthorized(w, r, "Unauthorized")
		return
	}

	var req struct {
		Count int `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, r, "Invalid request body")
		return
	}

	if req.Count <= 0 {
		response.WithValidationError(w, r, "Count must be greater than 0")
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
			response.WithInternalError(w, r, "Failed to generate upload URL")
			return
		}

		urls = append(urls, map[string]interface{}{
			"presignedUrl": url, // ← CHANGED: "upload_url" → "presignedUrl"
			"fileKey":      key, // ← CHANGED: "key" → "fileKey"
		})
	}

	// ← CHANGED: Return proper response structure
	response.WithPayload(w, r, map[string]interface{}{
		"presignedUrls": urls,
	})
}
