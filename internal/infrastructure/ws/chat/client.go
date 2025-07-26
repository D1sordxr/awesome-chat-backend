package chat

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	id       string
	chatIDs  map[string]bool
	conn     *websocket.Conn
	sendChan chan []byte
	close    atomic.Bool
	mu       sync.RWMutex
}

func NewClient(conn *websocket.Conn, id string, chatIDs ...string) *Client {
	chats := make(map[string]bool, len(chatIDs))
	for _, chatID := range chatIDs {
		chats[chatID] = true
	}

	client := &Client{
		id:       id,
		chatIDs:  chats,
		conn:     conn,
		sendChan: make(chan []byte, 256),
		close:    atomic.Bool{},
		mu:       sync.RWMutex{},
	}

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetPongHandler(func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	client.conn.SetCloseHandler(func(code int, text string) error {
		return client.Close()
	})

	return client
}

func (c *Client) Run(ctx context.Context) error {
	defer func() {
		c.close.Store(true)
		close(c.sendChan)
		_ = c.conn.Close()
	}()

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return c.readPump(ctx)
	})
	errGroup.Go(func() error {
		return c.writePump(ctx)
	})

	return errGroup.Wait()
}

func (c *Client) Close() error {
	if c.close.Swap(true) {
		return errors.New("already closed")
	}
	close(c.sendChan)
	return c.conn.Close()
}

func (c *Client) SendMessage(msg []byte) error {
	if c.close.Load() {
		return errors.New("client closed")
	}

	select {
	case c.sendChan <- msg:
		return nil
	default:
		return errors.New("client buffer overflow")
	}
}

func (c *Client) readPump(ctx context.Context) error {
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			_, _, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return fmt.Errorf("read error: %w", err)
				}
				return nil
			}
		}
	}
}

func (c *Client) writePump(ctx context.Context) error {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case message, ok := <-c.sendChan:
			if !ok {
				return c.conn.WriteMessage(websocket.CloseMessage, nil)
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return err
			}

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}

func (c *Client) JoinChat(chatID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.chatIDs[chatID] = true
}

func (c *Client) LeaveChat(chatID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.chatIDs, chatID)
}

func (c *Client) IsInChat(chatID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chatIDs[chatID]
}
