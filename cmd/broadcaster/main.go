package main

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/infrastructure/ws/chathub/consts"
	"awesome-chat/internal/infrastructure/ws/chathub/transport"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/rand"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	userID = "57c7ea1e-cedf-4ed8-bad2-ed9347baac70"
	chatID = "6585d0ad-f705-4723-8f9d-0b46c69290fa"
)

const (
	defaultOrigin  = "http://localhost/"
	serverURL      = "ws://localhost:8081/ws/" + userID
	messageRate    = time.Millisecond
	readBufferSize = 512
)

type WebSocketClient struct {
	conn    *websocket.Conn
	counter int
}

func main() {
	// Set up signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	// Initialize WebSocket connection
	client, err := NewWebSocketClient()
	if err != nil {
		log.Fatalf("Failed to initialize WebSocket client: %v", err)
	}
	defer client.Close()

	// Start message processing loop
	go client.ProcessMessages(done)

	// Wait for shutdown signal
	<-done
	log.Println("Shutting down client...")
}

func NewWebSocketClient() (*WebSocketClient, error) {
	conn, err := websocket.Dial(serverURL, "", defaultOrigin)
	if err != nil {
		return nil, fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	return &WebSocketClient{
		conn:    conn,
		counter: 0,
	}, nil
}

func (c *WebSocketClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}
}

func (c *WebSocketClient) ProcessMessages(done <-chan os.Signal) {
	ticker := time.NewTicker(messageRate)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := c.sendMessage(); err != nil {
				log.Printf("Error processing message: %v", err)
				return
			}
			c.counter++
		}
	}
}

func (c *WebSocketClient) sendMessage() error {
	// Prepare message
	message := dto.Message{
		UserID:  userID,
		ChatID:  chatID,
		Content: fmt.Sprintf("Hello World! %d", c.counter),
	}

	// Marshal message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Prepare operation DTO
	operation := transport.OperationHeader{
		ID:        rand.Int(),
		Operation: consts.SendMessage.String(),
		Body:      messageBytes,
	}

	// Marshal operation
	operationBytes, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	// Send message
	if _, err := c.conn.Write(operationBytes); err != nil {
		return fmt.Errorf("failed to write to WebSocket: %w", err)
	}

	// Read response
	response := make([]byte, readBufferSize)
	n, err := c.conn.Read(response)
	if err != nil {
		return fmt.Errorf("failed to read from WebSocket: %w", err)
	}

	log.Printf("Received response: %s", response[:n])
	return nil
}
