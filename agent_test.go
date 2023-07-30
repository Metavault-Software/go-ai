package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

func TestAgentStore_AddAgent(t *testing.T) {
	router := gin.Default()

	store := AgentStore{store: make(map[string]Agent)}
	storeSpec := AgentStoreSpec{}
	err := storeSpec.LoadAgentsFromFile("agents.json", &store)
	if err != nil {
		t.Errorf("Failed to load agents from file: %+v", err)
	}

	router.POST("/agents", store.AddAgent)
	router.GET("/agents", store.GetAgents)
	router.GET("/agents/:agentId", store.GetAgent)
	router.PUT("/agents/:agentId", store.UpdateAgent)
	router.DELETE("/agents/:agentId", store.DeleteAgent)

	err = router.Run()
	if err != nil {
		log.Printf("Failed to start server: %+v", err)
	} // Run on port 8080
}
