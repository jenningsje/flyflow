package requests

type StreamingData struct {
	ID               string         `json:"id"`
	Object           string         `json:"object"`
	Created          int64          `json:"created"`
	Model            string         `json:"model"`
	SystemFingerprint string        `json:"system_fingerprint"`
	Choices          []StreamingChoice `json:"choices"`
}

type StreamingChoice struct {
	Index        int         `json:"index"`
	Delta        StreamingDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

type StreamingDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
