package main

type Agent struct {
	AgentId  string                 `json:"id"`       // Unique identifier for each Agent
	Name     string                 `json:"name"`     // Job name for better understanding
	Executor string                 `json:"executor"` // Executor to be used for the agent
	Args     map[string]interface{} `json:"args"`     // Arguments for the executor
}
