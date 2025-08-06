package chathub

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	log ports.Logger

	id    string
	chats []string

	socket       *websocket.Conn
	send         chan []byte
	opChan       chan<- Operation
	respChanPool sync.Pool

	mu        sync.Mutex
	isClosed  atomic.Bool
	closeOnce sync.Once
}

func NewClient(
	log ports.Logger,
	socket *websocket.Conn,
	id string,
	opChan chan Operation,
	chats ...string,
) *Client {
	return &Client{
		log:    log,
		id:     id,
		socket: socket,
		send:   make(chan []byte, 256),
		opChan: opChan,
		chats:  chats,
		respChanPool: sync.Pool{
			New: func() interface{} {
				return make(chan OperationResponse, 1)
			},
		},
	}
}

func (c *Client) Run(ctx context.Context) error {
	c.log.Info("starting client session", "client_id", c.id)
	defer func() {
		_ = c.Close()
		c.log.Info("client session ended", "client_id", c.id)
	}()

	c.socket.SetPingHandler(nil)
	c.socket.SetPongHandler(func(string) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := c.socket.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			c.log.Error("failed to set read deadline", "client_id", c.id, "error", err)
		}
		return err
	})

	errGroup, gCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		c.log.Debug("starting write pump", "client_id", c.id)
		err := c.writePump(gCtx)
		if err != nil {
			c.log.Error("write pump failed", "client_id", c.id, "error", err)
		}
		return err
	})

	errGroup.Go(func() error {
		c.log.Debug("starting read pump", "client_id", c.id)
		err := c.readPump(gCtx)
		if err != nil {
			c.log.Error("read pump failed", "client_id", c.id, "error", err)
		}
		return err
	})

	return errGroup.Wait()
}

func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.isClosed.Store(true)
		c.log.Info("closing client connection", "client_id", c.id)

		c.mu.Lock()
		defer c.mu.Unlock()

		close(c.send)
		c.log.Debug("send channel closed", "client_id", c.id)

		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		err = c.socket.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(writeWait))
		if err != nil {
			c.log.Error("failed to send close message", "client_id", c.id, "error", err)
		}

		err = c.socket.Close()
		if err != nil {
			c.log.Error("failed to close socket", "client_id", c.id, "error", err)
		}
	})
	return err
}

func (c *Client) readPump(ctx context.Context) error {
	defer func() {
		_ = c.Close()
		c.log.Debug("read pump exiting", "client_id", c.id)
	}()

	c.mu.Lock()
	c.socket.SetReadLimit(maxMessageSize)
	_ = c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			c.log.Debug("context done in read pump", "client_id", c.id)
			return nil
		default:
			_, message, err := c.socket.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.log.Error("unexpected close error", "client_id", c.id, "error", err)
				} else {
					c.log.Debug("normal connection close", "client_id", c.id)
				}
				return err
			}

			c.log.Info("operation message received from client",
				"client_id", c.id,
				"content_length", len(message))

			if err = func() error {
				respChan := c.respChanPool.Get().(chan OperationResponse)
				defer func() {
					select {
					case <-respChan:
					default:
					}
					c.respChanPool.Put(respChan)
				}()

				opCtx, opCancel := context.WithTimeout(ctx, 5*time.Second)
				defer opCancel()

				op := Operation{
					ClientID: c.id,
					Data:     message,
					Retries:  0,
					RespChan: respChan,
					Ctx:      opCtx,
				}

				select {
				case c.opChan <- op:
				case <-ctx.Done():
					return nil
				case <-opCtx.Done():
					return opCtx.Err()
				}

				select {
				case resp := <-respChan:
					if opCtx.Err() != nil {
						return opCtx.Err()
					}
					c.send <- resp.ToJSON()
				case <-ctx.Done():
					return nil
				case <-opCtx.Done():
					return opCtx.Err()
				}

				return nil
			}(); err != nil {
				c.log.Error("operation error", "client_id", c.id, "error", err)
			}
		}
	}
}

func (c *Client) writePump(ctx context.Context) error {
	c.log.Debug("starting write pump", "client_id", c.id)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.log.Debug("write pump exiting", "client_id", c.id)
	}()

	for {
		select {
		case <-ctx.Done():
			c.log.Debug("context done in write pump", "client_id", c.id)
			return nil
		case message, ok := <-c.send:
			c.mu.Lock()
			if !ok {
				c.mu.Unlock()
				c.log.Debug("send channel closed, exiting write pump", "client_id", c.id)
				return nil
			}

			c.log.Info("sending message to client",
				"client_id", c.id,
				"content_length", len(message))

			err := c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.mu.Unlock()
				c.log.Error("failed to set write deadline", "client_id", c.id, "error", err)
				return err
			}

			err = c.socket.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				c.log.Error("failed to write message", "client_id", c.id, "error", err)
				return err
			}

			c.log.Debug("message successfully sent", "client_id", c.id)
		case <-ticker.C:
			c.mu.Lock()
			err := c.socket.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(writeWait),
			)
			c.mu.Unlock()

			if err != nil {
				c.log.Error("failed to send ping", "client_id", c.id, "error", err)
				return err
			}
			// c.log.Debug("ping sent", "client_id", c.id)
		}
	}
}

func LocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer func() { _ = conn.Close() }()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
