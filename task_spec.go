package main

// Define the executor interface and OpenAIAgent and GoogleTranslationAgent

type TaskSpec struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Executor string                 `json:"executor"`
	Args     map[string]interface{} `json:"args"`
}
