package repository

import (
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"gorm.io/gorm"
	"sync"
)

type DatabaseRepository struct {
	DB          *gorm.DB
	APIKeyCache sync.Map
	repo        Repository
}

func NewDatabaseRepository(db *gorm.DB, r Repository) *DatabaseRepository {
	repo := &DatabaseRepository{
		DB:   db,
		repo: r,
	}
	return repo
}

func (dr *DatabaseRepository) SaveQueryRecord(req *requests.CompletionRequest, resp *requests.CompletionResponse) error {
	jsonData, err := json.Marshal(req.Cr)
	if err != nil {
		return err
	}

	queryRecord := &models.QueryRecord{
		UserID:         0, // Set the appropriate user ID
		Request:        string(jsonData),
		Response:       resp.Response,
		RequestedModel: req.Cr.Model,
		MaxTokens:      req.Cr.MaxTokens,
		Stream:         req.Cr.Stream,
	}

	return dr.DB.Create(queryRecord).Error
}

func (dr *DatabaseRepository) ProxyRequest(r *requests.ProxyRequest) error {
	return dr.repo.ProxyRequest(r)
}

func (dr *DatabaseRepository) ChatCompletion(r *requests.CompletionRequest) (*requests.CompletionResponse, error) {
	resp, err := dr.repo.ChatCompletion(r)
	if err != nil {
		return nil, err
	}

	err = dr.SaveQueryRecord(r, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
