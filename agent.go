package main

type Agent struct {
	AgentId      string                 `json:"agent_id" firestore:"agent_id"` // Unique identifier for each Agent
	Name         string                 `json:"name" firestore:"name"`         // Job name for better understanding
	ExecutorType string                 `json:"executor" firestore:"executor"`
	Args         map[string]interface{} `json:"args" firestore:"args"` // Arguments for the executor
}
