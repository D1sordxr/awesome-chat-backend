package ws

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"context"
	"net/http"
)

type Broadcaster interface {
	Broadcast(ctx context.Context, message entity.OldMessage, payload []byte) error
}

type ConnManager interface {
	HandleWebSocket(
		w http.ResponseWriter,
		r *http.Request,
		header http.Header,
		userID string,
		initialChatIDs []string,
	) error
}

type Upgrader interface {
	HandleWebSocket(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		header http.Header,
		userID string,
		chatIDs ...string,
	) error
}

type MessageBroadcaster interface {
	Broadcast(ctx context.Context, message chathub.Message) error
}
