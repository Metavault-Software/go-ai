package main

type AgentType string

type Agent struct {
	AgentId  string // Unique identifier for each Agent
	Name     string // Job name for better understanding
	ArgTypes map[string]ValueType
}

type ValueType string

const (
	String  ValueType = "String"
	Int     ValueType = "Int"
	Float   ValueType = "Float"
	Boolean ValueType = "Boolean"
	Object  ValueType = "Object"
	Array   ValueType = "Array"
)
