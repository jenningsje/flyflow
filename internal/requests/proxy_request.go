package requests

import "net/http"

type ProxyRequest struct {
	R *http.Request
	W http.ResponseWriter
	APIKey string
	IsOpenAIKey bool
}
