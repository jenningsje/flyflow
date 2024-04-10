package webapp

import (
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/models"
	"net/http"
	"strconv"
	"time"
)

type TokensPerSecondData struct {
	Date          string  `json:"date"`
	TokensPerSecond float64 `json:"tokens_per_second"`
}

func (h *WebAppHandler) GetTokensPerSecondTimeSeries(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r, h.DB, h.Cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user's API keys
	var apiKeys []models.APIKey
	result := h.DB.Where("user_id = ?", userID).Find(&apiKeys)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve API keys", http.StatusInternalServerError)
		return
	}

	var plainAPIKeys []string
	for _, apiKey := range apiKeys {
		plainAPIKeys = append(plainAPIKeys, apiKey.Key)
	}

	// Get the aggregated query records for the past week
	lastWeek := time.Now().AddDate(0, 0, -7)
	var timeSeries []TokensPerSecondData
	result = h.DB.Table("query_records").
		Select("DATE_FORMAT(created_at, '%Y-%m-%dT%H') AS date, SUM((input_tokens + output_tokens) / request_time_seconds) AS tokens_per_second").
		Where("api_key IN (?) AND created_at >= ?", plainAPIKeys, lastWeek).
		Group("date").
		Order("date DESC").
		Find(&timeSeries)

	if result.Error != nil {
		http.Error(w, "Failed to retrieve query records", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(timeSeries)
}

func (h *WebAppHandler) GetQueryRecords(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r, h.DB, h.Cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user's API keys
	var apiKeys []models.APIKey
	result := h.DB.Where("user_id = ?", userID).Find(&apiKeys)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve API keys", http.StatusInternalServerError)
		return
	}

	var plainAPIKeys []string
	for _, apiKey := range apiKeys {
		plainAPIKeys = append(plainAPIKeys, apiKey.Key)
	}

	// Parse pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // Default page number
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default limit
	}

	// Get the paginated query records
	var queryRecords []models.QueryRecord
	result = h.DB.Where("api_key IN (?)", plainAPIKeys).Order("created_at desc").Limit(limit).Offset((page - 1) * limit).Find(&queryRecords)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve query records", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(queryRecords)
}