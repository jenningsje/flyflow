package repository

import (
	"github.com/flyflow-devs/flyflow/internal/logger"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
)

type AuthenticationRepository struct {
	DB          *gorm.DB
	APIKeyCache sync.Map
	repo        Repository
}

func NewAuthenticationRepository(db *gorm.DB, r Repository) *AuthenticationRepository {
	repo := &AuthenticationRepository{
		DB: db,
		repo: r,
	}
	go repo.startAPIKeyCacheUpdater(5 * time.Second)
	return repo
}

func (ar *AuthenticationRepository) IsValidAPIKey(apiKey string) bool {
	_, exists := ar.APIKeyCache.Load(apiKey)
	return exists
}

func (ar *AuthenticationRepository) ProxyRequest(r *requests.ProxyRequest) error {
	if r.IsOpenAIKey {
		return ar.repo.ProxyRequest(r)
	}
	if !ar.IsValidAPIKey(r.APIKey) {
		http.Error(r.W, "Invalid API Key", http.StatusUnauthorized)
	}

	// Passthrough the CompletionRequest without any authentication
	return ar.repo.ProxyRequest(r)
}

func (ar *AuthenticationRepository) ChatCompletion(r *requests.CompletionRequest) error {
	if r.IsOpenAIKey {
		return ar.repo.ChatCompletion(r)
	}
	if !ar.IsValidAPIKey(r.APIKey) {
		http.Error(r.W, "Invalid API Key", http.StatusUnauthorized)
	}

	// Passthrough the CompletionRequest without any authentication
	return ar.repo.ChatCompletion(r)
}

func (ar *AuthenticationRepository) startAPIKeyCacheUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ar.updateAPIKeyCache()
	}
}

func (ar *AuthenticationRepository) updateAPIKeyCache() {
	var apiKeys []models.APIKey
	result := ar.DB.Find(&apiKeys)
	if result.Error != nil {
		logger.S.Error("failed to updateAPIKeyCache ", result.Error )
	}

	newCache := sync.Map{}
	for _, key := range apiKeys {
		newCache.Store(key.Key, struct{}{})
	}

	ar.APIKeyCache = newCache
}