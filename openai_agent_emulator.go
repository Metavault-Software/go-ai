package main

import "context"

type EmulateOpenAIAgent struct {
	OpenAIAgent *OpenAIAgent
}

func (e EmulateOpenAIAgent) Execute(ctx context.Context, task *Task) error {
	return nil
}

func NewEmulateOpenAIAgent(spec TaskSpec) *EmulateOpenAIAgent {
	model := "gpt-3.5-turbo"
	agent := OpenAIAgent{}
	agent.Client = nil
	agent.Model = model
	messages := ToChatCompletionMessages(spec.Args)
	agent.Messages = messages
	return &EmulateOpenAIAgent{&agent}
}
