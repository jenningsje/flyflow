package requests

import (
	"github.com/flyflow-devs/flyflow/internal/models"
	"net/http"
)

type CompletionRequest struct {
	R *http.Request
	W http.ResponseWriter
	Cr InternalOpenAICompletionRequest
	APIKey string
	IsOpenAIKey bool
	Model *models.Model
}

type CompletionResponse struct {
	Response string
	ShouldSave bool
	RequestTimeSeconds float32
}

type OpenAICompletionRequest struct {
	Model       string      `json:"model,omitempty"`
	Messages    []Message   `json:"messages,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Tools       []Tool      `json:"tools,omitempty"`
	ToolChoice  string      `json:"tool_choice,omitempty"`
	LogProbs    bool        `json:"logprobs,omitempty"`
	TopLogProbs int         `json:"top_logprobs,omitempty"`
}

type InternalOpenAICompletionRequest struct {
	Model       string      `json:"model,omitempty"`
	Messages    []Message   `json:"messages,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Tools       []Tool      `json:"tools,omitempty"`
	ToolChoice  string      `json:"tool_choice,omitempty"`
	LogProbs    bool        `json:"logprobs,omitempty"`
	TopLogProbs int         `json:"top_logprobs,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Background  bool        `json:"background,omitempty"`
}

func (oar OpenAICompletionRequest) ToClaudeRequest() ClaudeRequest {
	cr := ClaudeRequest{
		Model:     oar.Model,
		MaxTokens: oar.MaxTokens,
		Metadata: &Metadata{
			Stream: oar.Stream,
		},
	}

	// Convert Messages to Claude's format
	// Initialize cr.Messages with the correct length
	cr.Messages = make([]ClaudeMessage, 0, len(oar.Messages))
	for _, msg := range oar.Messages {
		if msg.Role == "system" {
			// Cast msg.Content to string for the system message
			cr.System = msg.Content.(string)
		} else {
			// Append non-system messages to cr.Messages
			cr.Messages = append(cr.Messages, ClaudeMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	return cr
}

func (i InternalOpenAICompletionRequest) ToCompletionRequest() OpenAICompletionRequest {
	return OpenAICompletionRequest{
		Model: i.Model,
		Messages: i.Messages,
		MaxTokens: i.MaxTokens,
		Stream: i.Stream,
		Tools: i.Tools,
		ToolChoice: i.ToolChoice,
		LogProbs: i.LogProbs,
		TopLogProbs: i.TopLogProbs,
	}
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type Content struct {
	Type     string     `json:"type"`
	Text     string     `json:"text,omitempty"`
	ImageURL *ImageURL  `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}