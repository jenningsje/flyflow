package models

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	Name string
	Key string `gorm:"uniqueIndex"`
}