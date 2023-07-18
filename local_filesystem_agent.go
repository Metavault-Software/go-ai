package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type LocalFileSystemAgent struct {
	DirPath string
}

func NewLocalFileSystemAgent(taskSpec TaskSpec) *LocalFileSystemAgent {
	return &LocalFileSystemAgent{DirPath: taskSpec.Args["dir_path"].(string)}
}

func (a *LocalFileSystemAgent) Execute(ctx context.Context, task *Task) error {
	// Walk the directory and print file names
	err := filepath.Walk(a.DirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore directories
		if !info.IsDir() {
			fmt.Println(path)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	return nil
}
