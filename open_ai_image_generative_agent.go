package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OpenAIGenerativeAgent struct {
	Prompt string
	N      int
	Size   string
	Parser *CommandLineParser
}

func NewOpenAIGenerativeAgent(task Task) *OpenAIGenerativeAgent {
	agent := OpenAIGenerativeAgent{}
	agent.Parser = GetInstance()
	agent.Prompt = task.Args["prompt"].(string)
	agent.N = int(task.Args["n"].(float64))
	agent.Size = task.Args["size"].(string)
	return &agent
}

type GenerationRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func (a *OpenAIGenerativeAgent) Execute(ctx context.Context, task *Task) error {
	// Prepare the request body
	body := &GenerationRequest{
		Prompt: a.Prompt,
		N:      a.N,
		Size:   a.Size,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/images/generations",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	value := fmt.Sprintf("Bearer %s", a.Parser.OpenAPIKey)
	req.Header.Set("Authorization", value) // replace with your actual API key

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the HTTP response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response: %v", err)
	}

	// Here, respBytes will contain the generated image(s), which are typically returned as URLs or byte arrays
	// You would need to handle the response accordingly based on your application's needs
	fmt.Println(string(respBytes))

	return nil
}
