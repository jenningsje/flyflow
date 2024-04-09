package webapp

import (
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/google/uuid"
	"net/http"
)

func (h *WebAppHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r, h.DB, h.Cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var apiKeyReq struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&apiKeyReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	apiKey := models.APIKey{
		UserId: userID,
		Name:   apiKeyReq.Name,
		Key:    uuid.New().String(),
	}

	result := h.DB.Create(&apiKey)
	if result.Error != nil {
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(apiKey)
}

func (h *WebAppHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r, h.DB, h.Cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var apiKeys []models.APIKey
	result := h.DB.Where("user_id = ?", userID).Find(&apiKeys)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve API keys", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(apiKeys)
}

func (h *WebAppHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r, h.DB, h.Cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var apiKeyReq struct {
		Key string `json:"key"`
	}
	err = json.NewDecoder(r.Body).Decode(&apiKeyReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := h.DB.Where("user_id = ? AND key = ?", userID, apiKeyReq.Key).Delete(&models.APIKey{})
	if result.Error != nil {
		http.Error(w, "Failed to delete API key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
