package requests

import (
	"errors"
	"time"
)

type ClaudeRequest struct {
	Model    string          `json:"model"`
	Messages []ClaudeMessage `json:"messages"`
	System   string          `json:"system,omitempty"`
	MaxTokens int            `json:"max_tokens"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type ClaudeMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type Metadata struct {
	StopSequences []string `json:"stop_sequences,omitempty"`
	Stream        bool     `json:"stream,omitempty"`
	Temperature   float64  `json:"temperature,omitempty"`
	TopP          float64  `json:"top_p,omitempty"`
	TopK          int      `json:"top_k,omitempty"`
}

type ClaudeResponse struct {
	Content     []ContentBlock `json:"content"`
	ID          string         `json:"id"`
	Model       string         `json:"model"`
	Role        string         `json:"role"`
	StopReason  string         `json:"stop_reason"`
	StopSequence *string       `json:"stop_sequence"`
	Type        string         `json:"type"`
	Usage       Usage          `json:"usage"`
}

type ContentBlock struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func (cr *ClaudeResponse) ToOpenAIResponse() (*OpenAIResponse, error) {
	if len(cr.Content) == 0 {
		return nil, errors.New("internal server error")
	}
	return &OpenAIResponse{
		ID:               cr.ID,
		Object:           "chat.completion",
		Created:          time.Now().Unix(),
		Model:            cr.Model,
		SystemFingerprint: "",
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Message: OpenAIMessage{
					Role:    cr.Role,
					Content: cr.Content[0].Text,
				},
				Logprobs:     nil,
				FinishReason: cr.StopReason,
			},
		},
		Usage: OpenAIUsage{
			PromptTokens:     cr.Usage.InputTokens,
			CompletionTokens: cr.Usage.OutputTokens,
			TotalTokens:      cr.Usage.InputTokens + cr.Usage.OutputTokens,
		},
	}, nil
}

type OpenAIResponse struct {
	ID               string         `json:"id"`
	Object           string         `json:"object"`
	Created          int64          `json:"created"`
	Model            string         `json:"model"`
	SystemFingerprint string        `json:"system_fingerprint"`
	Choices          []OpenAIChoice `json:"choices"`
	Usage            OpenAIUsage    `json:"usage"`
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	Logprobs     interface{}   `json:"logprobs"`
	FinishReason string        `json:"finish_reason"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
