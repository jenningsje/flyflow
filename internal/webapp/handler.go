package webapp

import (
	"github.com/flyflow-devs/flyflow/internal/config"
	"gorm.io/gorm"
)

type WebAppHandler struct {
	DB *gorm.DB
	Cfg *config.Config
}

func NewWebAppHandler(db *gorm.DB, cfg *config.Config) *WebAppHandler {
	return &WebAppHandler{
		DB: db,
		Cfg: cfg,
	}
}
