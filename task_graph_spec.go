package main

import "encoding/json"

type TaskSpecGraph struct {
	Tasks []TaskSpec          `json:"tasks"`
	Edges map[string][]string `json:"edges"`
}

func (m TaskSpecGraph) Unmarshal(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), &m)
}
