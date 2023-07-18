package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type OpenAIAgent struct {
	Client   *openai.Client
	Model    string
	Messages []openai.ChatCompletionMessage
}

func NewOpenAIAgent(spec TaskSpec) *OpenAIAgent {
	parser := GetInstance()
	client := openai.NewClient(parser.OpenAPIKey)
	model := "gpt-3.5-turbo"
	agent := OpenAIAgent{}
	agent.Client = client
	agent.Model = model
	messages := ToChatCompletionMessages(spec.Args)
	agent.Messages = messages
	return &agent
}

func (e *OpenAIAgent) Execute(ctx context.Context, task *Task) error {
	req := openai.ChatCompletionRequest{
		Model:    e.Model,
		Messages: e.Messages,
	}
	resp, err := e.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		if err == context.Canceled {
			fmt.Println("OpenAI chat request was cancelled")
		} else {
			fmt.Println("OpenAI chat request failed with error:", err)
		}
		return err
	}
	fmt.Println("Chat response from OpenAI: ", resp.Choices[0].Message.Content)
	return nil
}

func ToChatCompletionMessages(args map[string]interface{}) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage

	for _, v := range args["messages"].([]interface{}) {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    v.(map[string]interface{})["role"].(string),
			Content: v.(map[string]interface{})["content"].(string),
		})
	}
	return messages
}
