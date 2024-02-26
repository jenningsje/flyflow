package repository

import "net/http"

type Repository interface {
	ProxyRequest(r *http.Request) (*http.Response, error)
}

