package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketServerAgent struct {
	Address   string
	clients   map[*websocket.Conn]bool // Registered clients
	broadcast chan string              // Broadcast channel
	mutex     sync.Mutex               // Mutex to protect clients
}

func NewWebSocketServerAgent(spec TaskSpec) *WebSocketServerAgent {
	return &WebSocketServerAgent{
		broadcast: make(chan string),
		clients:   make(map[*websocket.Conn]bool),
		Address:   spec.Args["address"].(string),
	}
}

func (agent *WebSocketServerAgent) Execute(ctx context.Context, task *Task) error {
	fmt.Printf("Starting WebSocket server agent on port 8080\n")
	http.HandleFunc("/ws", agent.handleConnections)
	go agent.handleMessages(ctx)
	go agent.Broadcast(ctx)
	err := http.ListenAndServe(agent.Address, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}

	return nil
}

func (agent *WebSocketServerAgent) handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade:", err)
		return
	}

	agent.mutex.Lock()
	agent.clients[ws] = true
	agent.mutex.Unlock()

	fmt.Println("New client connected")
}

func (agent *WebSocketServerAgent) handleMessages(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("Current time: %s", time.Now().String())
			agent.broadcast <- message
		case <-ctx.Done():
			fmt.Println("Stopping handleMessages routine")
			return
		}
	}
}

func (agent *WebSocketServerAgent) Broadcast(ctx context.Context) {
	for {
		select {
		case message := <-agent.broadcast:
			agent.mutex.Lock()
			for client := range agent.clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					client.Close()
					delete(agent.clients, client)
				}
			}
			agent.mutex.Unlock()
		case <-ctx.Done():
			fmt.Println("Stopping Broadcast routine")
			return
		}
	}
}
