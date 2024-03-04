package models

import "gorm.io/gorm"

type QueryRecord struct {
	gorm.Model
	QueryContext string
}
