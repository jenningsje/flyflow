package models

type QueryRecord struct {
	BaseModel
	APIKey            string    `json:"api_key" gorm:"index"`
	Request           string    `json:"request"`
	Response          string    `json:"response"`
	RequestedModel    string    `json:"requested_model"`
	MaxTokens         int       `json:"max_tokens"`
	InputTokens       int       `json:"input_tokens"`
	OutputTokens      int       `json:"output_tokens"`
	Temperature       float32   `json:"temperature"`
	TopP              float32   `json:"top_p"`
	PresencePenalty   float32   `json:"presence_penalty"`
	FrequencyPenalty  float32   `json:"frequency_penalty"`
	Stream            bool      `json:"stream"`
	Tags              []string  `json:"tags" gorm:"serializer:json"`
	RequestTimeSeconds float32  `json:"request_time_seconds"`
}