package chat

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/entity"
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Hub struct {
	clients   map[string]*Client
	chatRooms map[string]map[string]bool // chatID -> clientIDs
	mu        sync.RWMutex
	isClosed  atomic.Bool
	log       ports.Logger
	wg        sync.WaitGroup
}

func NewChatHub(log ports.Logger) *Hub {
	return &Hub{
		clients:   make(map[string]*Client),
		chatRooms: make(map[string]map[string]bool),
		log:       log,
	}
}

func (h *Hub) RegisterClient(client *Client, chatIDs []string) error {
	if h.isClosed.Load() {
		return errors.New("hub is closed")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client.id]; exists {
		return errors.New("client already registered")
	}

	h.clients[client.id] = client

	for _, chatID := range chatIDs {
		if _, exists := h.chatRooms[chatID]; !exists {
			h.chatRooms[chatID] = make(map[string]bool)
		}
		h.chatRooms[chatID][client.id] = true
	}

	h.wg.Add(1)
	return nil
}

func (h *Hub) UnregisterClient(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	defer h.wg.Done()

	client, exists := h.clients[clientID]
	if !exists {
		return
	}

	for chatID := range client.chatIDs {
		delete(h.chatRooms[chatID], clientID)
		if len(h.chatRooms[chatID]) == 0 {
			delete(h.chatRooms, chatID)
		}
	}

	delete(h.clients, clientID)
}

func (h *Hub) Broadcast(ctx context.Context, message entity.OldMessage, payload []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.isClosed.Load() {
		return errors.New("hub is closed")
	}

	clients, exists := h.chatRooms[message.ChatID]
	if !exists {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(clients))

	for clientID := range clients {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				if client, ok := h.clients[cid]; ok {
					if err := client.SendMessage(payload); err != nil {
						h.log.Error(
							"Error sending message to client",
							"chatID", message.ChatID,
							"clientID", cid,
							"error", err.Error(),
						)
						errChan <- err
					}
				}
			}
		}(clientID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) Start(_ context.Context) error {
	h.log.Info("Hub started")
	return nil
}

func (h *Hub) Shutdown(ctx context.Context) error {
	h.isClosed.Store(true)

	h.mu.Lock()
	for _, client := range h.clients {
		_ = client.Close()
	}
	h.mu.Unlock()

	done := make(chan struct{})
	go func() {
		h.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		h.log.Info("Hub shutdown complete")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
