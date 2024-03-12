package requests

import "net/http"

type CompletionRequest struct {
	R *http.Request
	W http.ResponseWriter
	Cr OpenAICompletionRequest
	APIKey string
	IsOpenAIKey bool
}

type CompletionResponse struct {
	Response string
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