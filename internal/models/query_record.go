package models

import "gorm.io/gorm"

type QueryRecord struct {
	gorm.Model
	UserID       uint
	Request      string
	Response     string
	RequestedModel        string
	MaxTokens    int
	Temperature  float32
	TopP         float32
	PresencePenalty float32
	FrequencyPenalty float32
	Stream       bool
}