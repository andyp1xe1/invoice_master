package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func llama(query string) Response {
	modelName := "tinyllama"

	requestBody, err := json.Marshal(Response{
		"model":  modelName,
		"prompt": query,
		"stream": false,
	})
	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
	}

	// Make the POST request
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	// Decode the response
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Response: %+v\n", result["response"])
	return result
}
