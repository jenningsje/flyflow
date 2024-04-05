package server

import (
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/config"
	"github.com/flyflow-devs/flyflow/internal/logger"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"github.com/flyflow-devs/flyflow/internal/webapp"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/flyflow-devs/flyflow/internal/repository"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Repo  repository.Repository
	DB    *gorm.DB
	Cfg   *config.Config
}

func NewServer(Config *config.Config, DB *gorm.DB,  repo repository.Repository) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		Repo:  repo,
		Cfg:   Config,
		DB:    DB,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.PathPrefix("/v1/chat/completion").HandlerFunc(s.handleCompletion)
	s.Router.PathPrefix("/").HandlerFunc(s.handleRequest)

	authHandler := webapp.NewAuthHandler(s.DB, s.Cfg)
	s.Router.HandleFunc("/webapp/signup", authHandler.SignUp).Methods(http.MethodPost)
	s.Router.HandleFunc("/webapp/login", authHandler.Login).Methods(http.MethodPost)
	s.Router.HandleFunc("/webapp/authcheck", authHandler.AuthCheck).Methods(http.MethodGet)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Extract the API key from the authentication header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	// Remove the "Bearer " prefix if present
	if len(apiKey) > 7 && strings.ToLower(apiKey[0:7]) == "bearer " {
		apiKey = apiKey[7:]
	}

	pr := &requests.ProxyRequest{
		R:      r,
		W: 		w,
		APIKey: apiKey,
	}

	if err := s.Repo.ProxyRequest(pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCompletion(w http.ResponseWriter, r *http.Request) {
	var req requests.InternalOpenAICompletionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the API key from the authentication header
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	// Remove the "Bearer " prefix if present
	if len(apiKey) > 7 && strings.ToLower(apiKey[0:7]) == "bearer " {
		apiKey = apiKey[7:]
	}

	cr := &requests.CompletionRequest{
		R:      r,
		W: 		w,
		Cr:     req,
		APIKey: apiKey,
	}

	if _, err := s.Repo.ChatCompletion(cr); err != nil {
		logger.S.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
