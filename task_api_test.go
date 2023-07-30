package main

import (
	"reflect"
	"testing"
)

func TestNewTaskStore(t *testing.T) {
	tests := []struct {
		name string
		want *TaskApi
	}{
		{
			name: "TestNewTaskStore",
			want: &TaskApi{Repo: NewFirestoreTaskRepository()},
		},
	}
	repo := NewFirestoreTaskRepository()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskApi(repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskApi() = %v, want %v", got, tt.want)
			}
		})
	}
}
