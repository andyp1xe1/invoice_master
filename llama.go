package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	modelName   = "llama-3.1-70b-versatile"
	modelTemp   = 1
	apiEndpoint = "https://api.groq.com/openai/v1/chat/completions"
)

var (
	systemPrompt string
	apiKey       string
)

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

func llama(docScan string) (Response, error) {
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

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Response: %+v\n", result)
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
