package chathub

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"

	"awesome-chat/internal/domain/app/ports"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	log ports.Logger

	chatClients map[string]map[string]*Client
	upgrader    *websocket.Upgrader

	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	errChan    chan error

	mu       sync.RWMutex
	wg       sync.WaitGroup
	isClosed atomic.Bool
}

func NewClientManager(log ports.Logger) *ClientManager {
	return &ClientManager{
		log:         log,
		chatClients: make(map[string]map[string]*Client),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		broadcast:  make(chan Message, 100),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		errChan:    make(chan error, 100),
	}
}

func (m *ClientManager) HandleWebSocket(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	header http.Header,
	userID string,
	chatIDs ...string,
) error {
	if m.isClosed.Load() {
		m.log.Warn("client manager is shutting down, rejecting new connection")
		return errors.New("client manager is shutting down")
	}

	m.log.Info("attempting to upgrade connection to WebSocket", "user_id", userID, "chat_ids", chatIDs)
	socket, err := m.upgrader.Upgrade(w, r, header)
	if err != nil {
		m.log.Error("WebSocket upgrade failed", "error", err.Error(), "user_id", userID)
		return err
	}

	client := NewClient(m.log, socket, userID, chatIDs...)
	m.log.Info("new client created", "client_id", client.id)

	select {
	case m.register <- client:
		m.log.Info("registering client in manager", "client_id", client.id)
		go func() {
			m.log.Debug("starting client goroutine", "client_id", client.id)
			if err = client.Run(context.WithoutCancel(ctx)); err != nil {
				m.log.Error("client run error", "client_id", client.id, "error", err.Error())
				m.errChan <- err
			}
			m.log.Info("client finished, unregistering", "client_id", client.id)
			m.unregister <- client
		}()
		return nil
	case <-ctx.Done():
		m.log.Warn("context canceled before client registration", "client_id", client.id)
		_ = socket.Close()
		return ctx.Err()
	}
}
func (m *ClientManager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, chat := range client.chats {
		if _, ok := m.chatClients[chat]; !ok {
			m.chatClients[chat] = make(map[string]*Client)
		}
		m.chatClients[chat][client.id] = client
		m.log.Debug("client registered successfully", "client_id", client.id, "chat_id", chat)
	}
}

func (m *ClientManager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, chat := range client.chats {
		if clients, ok := m.chatClients[chat]; ok {
			delete(clients, client.id)
			if len(clients) == 0 {
				delete(m.chatClients, chat)
			}
			m.log.Debug("client unregistered successfully", "client_id", client.id, "chat_id", chat)
		}
	}
}

func (m *ClientManager) broadcastToClients(message Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	msg := message.ToJSON()
	clients, ok := m.chatClients[message.ChatID]
	if !ok {
		m.log.Warn("no clients found for chat", "chat_id", message.ChatID)
		return
	}

	m.log.Info("broadcasting message to chat",
		"chat_id", message.ChatID,
		"sender_id", message.UserID,
		"content_length", len(message.Content))

	for id, client := range clients {
		select {
		case client.send <- msg:
			m.log.Debug("message queued for client",
				"client_id", id,
				"chat_id", message.ChatID,
				"queue_size", len(client.send))
		default:
			m.log.Warn("client buffer full, disconnecting",
				"client_id", id,
				"buffer_size", cap(client.send))
			go func(c *Client) {
				m.unregister <- c
				_ = c.Close()
			}(client)
		}
	}
}

func (m *ClientManager) Broadcast(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.broadcast <- message:
		m.log.Info("client broadcast", "message", message.ToJSON())
		return nil
	}
}

func (m *ClientManager) Start(ctx context.Context) error {
	m.wg.Add(1)
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-m.errChan:
			if !m.isClosed.Load() {
				m.log.Error("client error", "error", err.Error())
			}
		case client := <-m.register:
			m.registerClient(client)
		case client := <-m.unregister:
			m.unregisterClient(client)
			_ = client.Close()
		case message := <-m.broadcast:
			m.broadcastToClients(message)
		}
	}
}

func (m *ClientManager) Shutdown(ctx context.Context) error {
	if m.isClosed.Swap(true) {
		return nil
	}

	close(m.broadcast)
	close(m.register)
	close(m.unregister)
	close(m.errChan)

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
