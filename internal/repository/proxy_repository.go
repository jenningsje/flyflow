package repository

import (
	"net/http"

	"github.com/flyflow-devs/flyflow/internal/config"
)

type ProxyRepository struct{
	Config *config.Config
}

func NewProxyRepository(Config *config.Config) *ProxyRepository {
	return &ProxyRepository{
		Config: Config,
	}
}

func (pr *ProxyRepository) ProxyRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	r.URL.Host = "api.openai.com"
	r.URL.Scheme = "https"
	r.Header.Set("Authorization", "Bearer "+pr.Config.OpenAIAPIKey)
	return client.Do(r)
}
