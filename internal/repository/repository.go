package repository

import (
	"github.com/flyflow-devs/flyflow/internal/requests"
)

type Repository interface {
	ProxyRequest(r *requests.ProxyRequest) error
	ChatCompletion(r *requests.CompletionRequest) (*requests.CompletionResponse, error)
}
