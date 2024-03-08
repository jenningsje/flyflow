package repository

import (
	"bytes"
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/config"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"io"
	"net/http"
)

type ProxyRepository struct {
	Config *config.Config
}

func NewProxyRepository(Config *config.Config) *ProxyRepository {
	return &ProxyRepository{
		Config: Config,
	}
}

func (pr *ProxyRepository) ProxyRequest(r *requests.ProxyRequest) error {
	// Create a new request to avoid modifying the original request.
	newReq, err := http.NewRequest(r.R.Method, r.R.URL.String(), r.R.Body)
	if err != nil {
		return err
	}

	// Copy headers from the original request to the new request.
	for key, values := range r.R.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	// Set the new URL for the proxy request.
	newReq.URL.Host = "api.openai.com"
	newReq.URL.Scheme = "https"

	// Check if the user provided an OpenAI API key
	if r.IsOpenAIKey {
		newReq.Header.Set("Authorization", "Bearer "+r.APIKey)
	} else {
		newReq.Header.Set("Authorization", "Bearer "+pr.Config.OpenAIAPIKey)
	}

	// Use a new HTTP client to make the request.
	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the response headers to the response writer
	for key, values := range resp.Header {
		for _, value := range values {
			r.W.Header().Add(key, value)
		}
	}

	// Write the response status code
	r.W.WriteHeader(resp.StatusCode)

	// Copy the response body to the response writer
	_, err = io.Copy(r.W, resp.Body)
	return err
}

func (pr *ProxyRepository) ChatCompletion(r *requests.CompletionRequest) error {
	jsonData, err := json.Marshal(r.Cr)
	if err != nil {
		return err
	}

	// Create a new HTTP request with the JSON data
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the content type and authorization headers
	req.Header.Set("Content-Type", "application/json")
	if r.IsOpenAIKey {
		req.Header.Set("Authorization", "Bearer "+r.APIKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+pr.Config.OpenAIAPIKey)
	}

	// Use a new HTTP client to make the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the response headers to the response writer
	for key, values := range resp.Header {
		for _, value := range values {
			r.W.Header().Add(key, value)
		}
	}

	// Write the response status code
	r.W.WriteHeader(resp.StatusCode)

	// Copy the response body to the response writer
	_, err = io.Copy(r.W, resp.Body)
	return err
}