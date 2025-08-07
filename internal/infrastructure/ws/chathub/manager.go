package chathub

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/ws/chathub/consts"
	chathubErrors "awesome-chat/internal/infrastructure/ws/chathub/errors"
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
)

type ClientManager struct {
	log ports.Logger

	chatClients map[string]map[string]*Client
	upgrader    *websocket.Upgrader

	opChan     chan Operation
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	errChan    chan error

	opHandler operationHandler

	mu       sync.RWMutex
	wg       sync.WaitGroup
	isClosed atomic.Bool
}

func NewClientManager(
	log ports.Logger,
) *ClientManager {
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
		opChan:     make(chan Operation, consts.ChanBuff),
		broadcast:  make(chan Message, consts.ChanBuff),
		register:   make(chan *Client, consts.ChanBuff),
		unregister: make(chan *Client, consts.ChanBuff),
		errChan:    make(chan error, consts.ChanBuff),
	}
}

func (m *ClientManager) MustOperationHandler(handler operationHandler) {
	m.opHandler = handler
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

	client := NewClient(m.log, socket, userID, m.opChan, chatIDs...)
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
		if client.isClosed.Load() {
			m.log.Debug("skipping message for client",
				"client_id", id,
				"chat_id", message.ChatID,
				"queue_size", len(client.send),
				"details", "client is closed",
			)
			continue
		}
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
		case op := <-m.opChan:
			opResp := m.opHandler.Handle(op.Ctx, op.Data)
			if opResp.Error != nil {
				m.log.Error("operation error",
					"client_id", op.ClientID, "error", opResp.Error.Error(), "retries", op.Retries,
				)
				switch {
				case errors.Is(opResp.Error, context.Canceled):
					op.RespChan <- opResp
				case errors.Is(opResp.Error, context.DeadlineExceeded):
					op.RespChan <- opResp
				case errors.Is(opResp.Error, chathubErrors.ErrUnsupportedOp):
					op.RespChan <- opResp
				case errors.Is(opResp.Error, chathubErrors.ErrInvalidOpFormat):
					op.RespChan <- opResp
				default:
					if op.Retries < 3 {
						op.Retries++
						select {
						case m.opChan <- op:
							m.log.Debug("retrying operation",
								"client_id", op.ClientID, "retries", op.Retries,
							)
							continue
						case <-ctx.Done():
							return ctx.Err()
						}
					} else {
						m.log.Error("operation ",
							"operation_details", op, "last_error", opResp.Error.Error(),
						)
						op.RespChan <- opResp
					}
				}
			}
		}
	}
}

func (m *ClientManager) Shutdown(ctx context.Context) error {
	if m.isClosed.Swap(true) {
		return nil
	}

	close(m.register)
	close(m.unregister)
	close(m.opChan)
	close(m.broadcast)
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
