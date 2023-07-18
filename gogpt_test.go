package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"testing"
)

func TestTask_Run(t *testing.T) {
	parser := CommandLineParser{}

	// Read input file
	inputData, err := ioutil.ReadFile(parser.InputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	// Create OpenAI client
	client := openai.NewClient(parser.OpenAPIKey)

	// Define the chat message content
	content := string(inputData)

	// Determine the model
	var model string
	switch parser.Model {
	case "gpt3dot5":
		model = openai.GPT3Dot5Turbo
	case "gpt4":
		model = openai.GPT4
	default:
		log.Fatalf("Invalid model specified. Please choose either 'gpt3dot5' or 'gpt4'.")
	}

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}

	// Send chat completion request
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Chat completion error: %v\n", err)
	}

	// Extract the generated completion from the response
	completion := resp.Choices[0].Message.Content

	// Write completion to the output file
	err = ioutil.WriteFile(parser.OutputFile, []byte(completion), 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Output written to: %s\n", parser.OutputFile)
}
