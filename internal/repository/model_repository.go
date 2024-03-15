package repository

import (
	"errors"
	"github.com/flyflow-devs/flyflow/internal/logger"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ModelRepository struct {
	DB          *gorm.DB
	ModelCache sync.Map
	repo        Repository
}

func NewModelRepository(db *gorm.DB, r Repository) *ModelRepository {
	repo := &ModelRepository{
		DB:   db,
		repo: r,
	}
	go repo.startModelCacheUpdater(5 * time.Second)
	return repo
}

func (ar *ModelRepository) ProxyRequest(r *requests.ProxyRequest) error {
	// Passthrough the CompletionRequest without any authentication
	return ar.repo.ProxyRequest(r)
}

func (ar *ModelRepository) ChatCompletion(r *requests.CompletionRequest) (*requests.CompletionResponse, error) {
	// Get the model name from the OpenAI completion request
	modelName := r.Cr.Model

	// Retrieve the model from the cache based on the model name
	modelValue, ok := ar.ModelCache.Load(modelName)
	if !ok {
		return nil, errors.New("model not found in cache")
	}

	model, ok := modelValue.(*models.Model)
	if !ok {
		return nil, errors.New("invalid model type in cache")
	}

	// Set the retrieved model on the completion request
	r.Model = model
	r.Cr.Model = model.InternalModelName

	// Call the ChatCompletion method of the underlying repository
	return ar.repo.ChatCompletion(r)
}

func (ar *ModelRepository) startModelCacheUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ar.updateModelCache()
	}
}

func (ar *ModelRepository) updateModelCache() {
	var llmModels []models.Model
	result := ar.DB.Find(&llmModels)
	if result.Error != nil {
		logger.S.Error("failed to updateModelCache ", result.Error)
	}

	newCache := sync.Map{}
	for _, m := range llmModels {
		newCache.Store(m.ModelName, &m)
	}

	ar.ModelCache = newCache
}