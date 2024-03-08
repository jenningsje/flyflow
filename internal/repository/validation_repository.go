package repository

import (
	"github.com/flyflow-devs/flyflow/internal/requests"
	"regexp"
)

type ValidationRepository struct{
	repo Repository
}

func NewValidationRepository(repo Repository) *ValidationRepository {
	return &ValidationRepository{
		repo: repo,
	}
}

func (vr *ValidationRepository) ProxyRequest(r *requests.ProxyRequest) error {
	if match, _ := regexp.MatchString(`^sk-[a-zA-Z0-9]{48}$`, r.APIKey); match {
		r.IsOpenAIKey = true
	} else {
		r.IsOpenAIKey = false
	}
	return vr.repo.ProxyRequest(r)
}

func (vr *ValidationRepository) ChatCompletion(r *requests.CompletionRequest) error {
	if match, _ := regexp.MatchString(`^sk-[a-zA-Z0-9]{48}$`, r.APIKey); match {
		r.IsOpenAIKey = true
	} else {
		r.IsOpenAIKey = false
	}

	return vr.repo.ChatCompletion(r)
}