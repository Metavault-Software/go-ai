package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestTaskDag(t *testing.T) {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Open JSON file
	jsonFile, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalf("Error closing JSON file: %v", err)
		}
	}(jsonFile)

	// Read JSON file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Parse JSON into task specifications
	var taskDagSpec TaskDagSpec
	err = json.Unmarshal(byteValue, &taskDagSpec)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Create the TaskDag from the TaskDagSpec
	taskDag := new(TaskDag)

	err = taskDag.FromTaskDagSpec(taskDagSpec)
	if err != nil {
		log.Fatalf("Error creating TaskDag: %v", err)
	}

	started := make(chan Task)
	result := make(chan Task)
	go taskDag.RunConcurrently(started, result)

	for {
		select {
		case task, ok := <-started:
			if ok {
				fmt.Printf("Started: %s\n", task.Id)
			}
		case task, ok := <-result:
			if ok {
				fmt.Printf("Finished: %s\n", task.Id)
			} else {
				// No more tasks to process
				return
			}
		default:
			// Let's not hog the CPU
			time.Sleep(50 * time.Millisecond)
		}
	}
}
