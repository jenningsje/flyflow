package models

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	UserId uint
	ModelName string
	InternalModelName string
	APIUrl string
	APIKey string
}