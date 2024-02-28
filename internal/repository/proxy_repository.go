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
	// Create a new request to avoid modifying the original request.
	newReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers from the original request to the new request.
	for key, values := range r.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	// Set the new URL for the proxy request.
	newReq.URL.Host = "api.openai.com"
	newReq.URL.Scheme = "https"

	// Set the Authorization header with the API key.
	newReq.Header.Set("Authorization", "Bearer "+pr.Config.OpenAIAPIKey)

	// Use a new HTTP client to make the request.
	client := &http.Client{}
	return client.Do(newReq)
}

