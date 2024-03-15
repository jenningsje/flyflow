package models

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	UserId uint `gorm:"index"`
	Name string
	Key string `gorm:"uniqueIndex"`
}