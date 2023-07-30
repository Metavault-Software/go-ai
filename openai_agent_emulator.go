package main

import "context"

type EmulateOpenAIAgent struct {
	OpenAIAgent *OpenAIAgent
}

func (e EmulateOpenAIAgent) Execute(ctx context.Context, task *Task) error {
	return nil
}

func NewEmulateOpenAIAgent(task Task) *EmulateOpenAIAgent {
	model := "gpt-3.5-turbo"
	agent := OpenAIAgent{}
	agent.Client = nil
	agent.Model = model
	messages := ToChatCompletionMessages(task.Args)
	agent.Messages = messages
	return &EmulateOpenAIAgent{&agent}
}
