package main


type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint *string  `json:"system_fingerprint"` // optional field
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

// Choice represents each choice in the response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Logprobs     *string `json:"logprobs"` // optional field
}

// Message represents the assistant's message content
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage represents the usage statistics of the completion
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTime        float64 `json:"total_time"`
}
