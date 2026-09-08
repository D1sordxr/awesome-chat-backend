package chat

import (
	"context"

	"awesome-chat/internal/application/chat/command"
	"awesome-chat/internal/application/chat/query"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/handler_mocks.go -package=mocks

type useCase interface {
	Create(ctx context.Context, in command.CreateIn) (command.CreateOut, error)
	AddMember(ctx context.Context, in command.AddMemberIn) (command.AddMemberOut, error)
	Previews(ctx context.Context, in query.PreviewsIn) (query.PreviewsOut, error)
}
