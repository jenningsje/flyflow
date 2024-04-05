package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex"`
	HashedPassword string
	Password       string `gorm:"-"` // Exclude from database
}