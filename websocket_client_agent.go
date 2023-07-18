package main

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WebSocketClientAgent struct {
	Addr string
}

func NewWebSocketClientAgent(spec TaskSpec) *WebSocketClientAgent {
	addr := spec.Args["address"].(string)
	return &WebSocketClientAgent{Addr: addr}
}

func (e *WebSocketClientAgent) Execute(ctx context.Context, task *Task) error {
	time.Sleep(5 * time.Second)
	c, _, err := websocket.DefaultDialer.Dial(e.Addr, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("Received: %s", message)
			}
		}
	}()

	<-done

	return nil
}
