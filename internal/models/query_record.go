package models

import "gorm.io/gorm"

type QueryRecord struct {
	gorm.Model
	APIKey       string `gorm:"index"`
	Request      string
	Response     string
	RequestedModel        string
	MaxTokens    int
	InputTokens  int
	OutputTokens int
	Temperature  float32
	TopP         float32
	PresencePenalty float32
	FrequencyPenalty float32
	Stream       bool
	Tags         []string `gorm:"serializer:json"`
}