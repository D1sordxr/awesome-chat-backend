package chathub

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/ws/chathub/consts"
	chathubErrors "awesome-chat/internal/infrastructure/ws/chathub/errors"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
)

type ClientManagerV2 struct {
	log ports.Logger

	clientStore ClientStore
	upgrader    *websocket.Upgrader

	broadcast chan Message
	opChan    chan Operation
	errChan   chan error

	opHandler operationHandler

	mu       sync.RWMutex
	wg       sync.WaitGroup
	isClosed atomic.Bool
}

func NewClientManagerV2(
	log ports.Logger,
	clientStore ClientStore,
) *ClientManagerV2 {
	return &ClientManagerV2{
		log:         log,
		clientStore: clientStore,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		opChan:    make(chan Operation, consts.ChanBuff),
		broadcast: make(chan Message, consts.ChanBuff),
		errChan:   make(chan error, consts.ChanBuff),
	}
}

func (m *ClientManagerV2) MustSetOperationHandler(handler operationHandler) {
	m.opHandler = handler
}

func (m *ClientManagerV2) HandleWebSocket(
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

	go func() {
		m.clientStore.Add(client)
		defer func() {
			if rr := recover(); rr != nil {
				m.errChan <- fmt.Errorf("panic: %v", rr)
			}
			if clientErr := client.Close(); clientErr != nil {
				m.errChan <- clientErr
			}
			m.clientStore.Remove(client)
		}()

		if err = client.Run(context.WithoutCancel(ctx)); err != nil {
			m.errChan <- fmt.Errorf("client error: %w", err)
		}
	}()

	return nil
}

func (m *ClientManagerV2) broadcastToClients(message Message) {
	opResp := &OperationResponse{
		OperationType: consts.Broadcast.String(),
		Success:       true,
		Data:          message,
	}
	opRespBytes := opResp.ToJSON()

	clients, ok := m.clientStore.GetClients(message.ChatID)
	if !ok {
		m.log.Warn("no clients found for chat", "chat_id", message.ChatID)
		return
	}

	m.log.Info("broadcasting message to chat",
		"chat_id", message.ChatID,
		"sender_id", message.UserID,
		"content_length", len(message.Content),
	)

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
		case client.send <- opRespBytes:
			m.log.Debug("message queued for client",
				"client_id", id,
				"chat_id", message.ChatID,
				"queue_size", len(client.send))
		default:
			m.log.Warn("client buffer full, disconnecting",
				"client_id", id,
				"buffer_size", cap(client.send))
			go func(c *Client) {
				_ = c.Close()
			}(client)
		}
	}
}

func (m *ClientManagerV2) Start(ctx context.Context) error {
	m.wg.Add(1)
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-m.errChan:
			if !m.isClosed.Load() {
				m.log.Error("ClientManager received error", "error", err.Error())
			}
		case message := <-m.broadcast:
			if !m.isClosed.Load() {
				m.broadcastToClients(message)
			}
		case op := <-m.opChan:
			if m.isClosed.Load() {
				op.RespChan <- OperationResponse{Error: errors.New("client manager is shutting down")}
				continue
			}
			opResp := m.opHandler.Handle(op.Ctx, op.Data)
			if opResp.Error == nil {
				op.RespChan <- opResp
			} else {
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

func (m *ClientManagerV2) Broadcast(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.broadcast <- message:
		m.log.Info("client broadcast", "chat_id", message.ChatID)
		return nil
	}
}

func (m *ClientManagerV2) Shutdown(ctx context.Context) error {
	if m.isClosed.Swap(true) {
		return nil
	}

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
