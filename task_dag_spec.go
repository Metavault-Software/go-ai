package main

import "encoding/json"

type TaskDagSpec struct {
	Tasks []TaskSpec          `json:"tasks"`
	Edges map[string][]string `json:"edges"`
}

func (m TaskDagSpec) Unmarshal(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), &m)
}
