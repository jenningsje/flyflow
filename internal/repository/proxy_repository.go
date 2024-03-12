package repository

import (
	"bytes"
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/config"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"io"
	"net/http"
	"strings"
)

type ProxyRepository struct {
	Config *config.Config
	Client *http.Client
}

func NewProxyRepository(Config *config.Config) *ProxyRepository {
	client := &http.Client{}

	return &ProxyRepository{
		Config: Config,
		Client: client,
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

	var apiKey string
	if r.IsOpenAIKey {
		apiKey = r.APIKey
	} else {
		apiKey = pr.Config.OpenAIAPIKey
	}

	// Set the Authorization header with the API key.
	newReq.Header.Set("Authorization", "Bearer "+apiKey)

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

func (pr *ProxyRepository) ChatCompletion(r *requests.CompletionRequest) (*requests.CompletionResponse, error) {
	jsonData, err := json.Marshal(r.Cr)
	if err != nil {
		return &requests.CompletionResponse{}, err
	}

	// Create a new HTTP request with the JSON data
	req, err := http.NewRequest(r.R.Method, r.R.URL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return &requests.CompletionResponse{}, err
	}

	req.URL.Host = "api.openai.com"
	req.URL.Scheme = "https"

	// Set the content type and authorization headers
	req.Header.Set("Content-Type", "application/json")
	if r.IsOpenAIKey {
		req.Header.Set("Authorization", "Bearer "+r.APIKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+pr.Config.OpenAIAPIKey)
	}

	// Use a new HTTP client to make the request.
	resp, err := pr.Client.Do(req)
	if err != nil {
		return &requests.CompletionResponse{}, err
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
	// Create a buffer to store the response data
	buffer := make([]byte, 1024)

	// Stream the response body to the response writer
	var responseBuilder strings.Builder
	for {
		// Read the response data into the buffer
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Write the data to the response writer
		_, err = r.W.Write(buffer[:n])
		if err != nil {
			return nil, err
		}

		// Append the response data to the responseBuilder
		responseBuilder.Write(buffer[:n])

		// Flush the response writer if it implements the http.Flusher interface
		if flusher, ok := r.W.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	return &requests.CompletionResponse{
		Response: responseBuilder.String(),
	}, nil
}
