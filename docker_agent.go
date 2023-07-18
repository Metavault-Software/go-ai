package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"os"
	_ "path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type DockerAgent struct {
	Image   string
	Command []string
}

func NewDockerAgent(spec TaskSpec) Executor {
	agent := DockerAgent{}
	agent.Image = spec.Args["image"].(string)
	agent.Command = agent.FromArgs(spec.Args["command"].([]interface{}))
	return &agent
}

func (a *DockerAgent) Execute(ctx context.Context, task *Task) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: a.Image,
			Cmd:   a.Command,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: cwd,
					Target: "/workdir",
				},
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the container to finish
	_, errors := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	// drain the error chan and return the first error if any
	err = drainErrors(errors)
	return err
}

func (a *DockerAgent) FromArgs(i []interface{}) []string {
	s := make([]string, len(i))
	for i, v := range i {
		s[i] = v.(string)
	}
	return s
}

func drainErrors(errors <-chan error) error {
	var firstError error

	// Drain the error chan
	for err := range errors {
		// Capture the first non-nil error
		if firstError == nil && err != nil {
			firstError = err
		}
	}

	// Return the first non-nil error, if any
	return firstError
}
