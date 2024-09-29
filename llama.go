package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Define the structures according to the API response
//type ChatCompletion struct {
//	ID      string   `json:"id"`
//	Object  string   `json:"object"`
//	Created int64    `json:"created"`
//	Model   string   `json:"model"`
//	Choices []Choice `json:"choices"`
//	Usage   Usage    `json:"usage"`
//}
//
//type Choice struct {
//	Index        int     `json:"index"`
//	Message      Message `json:"message"`
//	FinishReason string  `json:"finish_reason"`
//	Logprobs     *string `json:"logprobs"` // optional field
//}
//
//type Message struct {
//	Role    string `json:"role"`
//	Content string `json:"content"`
//}
//
//type Usage struct {
//	PromptTokens     int `json:"prompt_tokens"`
//	CompletionTokens int `json:"completion_tokens"`
//	TotalTokens      int `json:"total_tokens"`
//}

type ChatCompletion struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
	modelName   = "llama3.2"
	modelTemp   = 0.7
)

var (
	apiKey       = os.Getenv("")
	systemPrompt string
)

func llama(docScan string) (*ChatCompletion, error) {
	requestBody, err := json.Marshal(Response{
		"model":  modelName,
		"system": systemPrompt,
		"prompt": docScan,
		"stream": false,
	})

	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
		return nil, err
	}

	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d", resp.StatusCode)
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var result ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
		return nil, err
	}

	fmt.Printf("Response: %+v\n", result)
	return &result, nil
}

//func llama(docScan string) (*ChatCompletion, error) {
//	headers := map[string]string{
//		"Content-Type":  "application/json",
//		"Authorization": "Bearer " + apiKey,
//	}
//
//	requestBody, err := json.Marshal(map[string]interface{}{
//		"messages": []Message{
//			{
//				Role:    "system",
//				Content: systemPrompt,
//			},
//			{
//				Role:    "user",
//				Content: docScan,
//			},
//		},
//		"model":       modelName,
//		"temperature": modelTemp,
//		"max_tokens":  1024,
//		"top_p":       modelTemp,
//		"stream":      false,
//	})
//
//	if err != nil {
//		log.Fatalf("Error marshaling request body: %v", err)
//		return nil, err
//	}
//
//	resp, err := postWithHeaders(apiEndpoint, requestBody, headers)
//	defer resp.Body.Close()
//
//	if err != nil {
//		log.Printf("Error making POST request: %v", err)
//		return nil, err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		log.Printf("Error: received status code %d", resp.StatusCode)
//		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
//	}
//
//	var result ChatCompletion
//	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//		log.Fatalf("Error decoding response: %v", err)
//		return nil, err
//	}
//
//	fmt.Printf("Response: %+v\n", result)
//	return &result, nil
//}
//
//func postWithHeaders(url string, jsonData []byte, headers map[string]string) (*http.Response, error) {
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
//	if err != nil {
//		return nil, err
//	}
//
//	for key, value := range headers {
//		req.Header.Set(key, value)
//	}
//
//	client := &http.Client{}
//	return client.Do(req)
//}

func readSys(path string) (string, error) {
	fd, err := os.Open(path)
	defer fd.Close()
	if err != nil {
		return "", err
	}
	buff := make([]byte, 1024)
	fd.Read(buff)
	return string(buff), nil
}

func llama(docScan string) (*ChatCompletion, error) {
	slog.Info("Key xd", apiKey)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	requestBody, err := json.Marshal(
		Response{
			"messages": []Response{
				Response{
					"role":    "system",
					"content": systemPrompt,
				},
				{
					"role":    "user",
					"content": docScan,
				},
			},
			"model":       modelName,
			"temperature": 1,
			"max_tokens":  1024,
			"top_p":       modelTemp,
			"stream":      false,
			"response_format": Response{
				"type": "json_object",
			},
			"stop": nil,
		},
	)

	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
	}

	resp, err := postWithHeaders(apiEndpoint, requestBody, headers)
	defer resp.Body.Close()

	if err != nil {
		slog.Error("Error making POST request: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error: received status code %d", resp.StatusCode)
	}

	var result *ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	return result, nil
}

func postWithHeaders(url string, jsonData []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}
