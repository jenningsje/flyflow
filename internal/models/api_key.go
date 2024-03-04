package models

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	Key string `gorm:"uniqueIndex"`
}