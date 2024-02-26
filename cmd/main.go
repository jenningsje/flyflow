package main

import (
	"log"
	"net/http"

	"github.com/flyflow-devs/flyflow/internal/config"
	"github.com/flyflow-devs/flyflow/internal/server"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize server
	s := server.NewServer(cfg)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}
