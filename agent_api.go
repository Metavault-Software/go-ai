package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AgentStoreSpec struct {
	Agents []Agent `json:"agents"`
}

type AgentStore struct {
	store map[string]Agent
}

// LoadAgentsFromFile Method to load agents from a JSON file
func (ss *AgentStoreSpec) LoadAgentsFromFile(filename string, s *AgentStore) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Printf("Failed to close file: %+v", err)
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	// Unmarshal JSON data

	err = json.Unmarshal(byteValue, ss)
	if err != nil {
		return err
	}

	// Add agents to store
	for _, agent := range ss.Agents {
		s.store[agent.AgentId] = agent
	}

	return nil
}

func (s *AgentStore) AddAgent(c *gin.Context) {
	var newAgent Agent
	if err := c.BindJSON(&newAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.store[newAgent.AgentId] = newAgent
	c.JSON(http.StatusOK, newAgent)
}

func (s *AgentStore) GetAgents(c *gin.Context) {
	agents := make([]Agent, 0, len(s.store))
	for _, agent := range s.store {
		agents = append(agents, agent)
	}
	c.JSON(http.StatusOK, agents)
}

func (s *AgentStore) GetAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent, ok := s.store[agentId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}
	c.JSON(http.StatusOK, agent)
}

func (s *AgentStore) UpdateAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	var updatedAgent Agent
	if err := c.BindJSON(&updatedAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.store[agentId] = updatedAgent
	c.JSON(http.StatusOK, updatedAgent)
}

func (s *AgentStore) DeleteAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	delete(s.store, agentId)
	c.JSON(http.StatusOK, gin.H{"status": "Agent deleted"})
}
