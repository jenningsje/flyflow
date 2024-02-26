package server

import (
	"github.com/flyflow-devs/flyflow/internal/config"
	"io"
	"net/http"

	"github.com/flyflow-devs/flyflow/internal/repository"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Repo  repository.Repository
}

func NewServer(Config *config.Config) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		Repo:  repository.NewProxyRepository(Config),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.PathPrefix("/").HandlerFunc(s.handleRequest)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := s.Repo.ProxyRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy body
	io.Copy(w, resp.Body)
}
