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
		Select("to_char(created_at, 'YYYY-MM-DD\"T\"HH24') AS date, AVG(output_tokens / request_time_seconds) AS tokens_per_second").
		Where("api_key IN (?) AND created_at >= ? AND request_time_seconds > 0", plainAPIKeys, lastWeek).
		Group("date").
		Order("date ASC").
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
	createdAtStr := r.URL.Query().Get("created_at")
	limitStr := r.URL.Query().Get("limit")

	var createdAt time.Time
	if createdAtStr != "" {
		createdAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			http.Error(w, "Invalid created_at value", http.StatusBadRequest)
			return
		}
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default limit
	}

	// Get the paginated query records
	var queryRecords []models.QueryRecord
	query := h.DB.Where("api_key IN (?)", plainAPIKeys).Order("created_at desc").Limit(limit)

	if !createdAt.IsZero() {
		query = query.Where("created_at < ?", createdAt)
	}

	result = query.Find(&queryRecords)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve query records", http.StatusInternalServerError)
		return
	}

	response := struct {
		Data       []models.QueryRecord `json:"data"`
		CreatedAt  time.Time            `json:"created_at"`
		HasMore    bool                 `json:"has_more"`
	}{
		Data:      queryRecords,
		CreatedAt: createdAt,
		HasMore:   len(queryRecords) == limit,
	}

	json.NewEncoder(w).Encode(response)
}