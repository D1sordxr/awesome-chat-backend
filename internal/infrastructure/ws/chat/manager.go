package chat

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type ConnManager struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewManager(hub *Hub) *ConnManager {
	return &ConnManager{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { // TODO
				return true
			},
		},
	}
}

func (m *ConnManager) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
	header http.Header,
	userID string,
	initialChatIDs []string,
) error {
	if m.hub.isClosed.Load() {
		return errors.New("hub is closed")
	}

	header = http.Header{}
	header.Set("Sec-WebSocket-Protocol", "chat")

	conn, err := m.upgrader.Upgrade(w, r, header)
	if err != nil {
		return err
	}

	client := NewClient(conn, userID, initialChatIDs...)

	if err = m.hub.RegisterClient(client, initialChatIDs); err != nil {
		_ = conn.Close()
		return err
	}

	go func() {
		defer func() {
			_ = client.Close()
			m.hub.UnregisterClient(userID)
		}()

		if err = client.Run(r.Context()); err != nil {
			m.hub.log.Error("Client error", "clientID", userID, "error", err)
		}
	}()

	return nil
}
