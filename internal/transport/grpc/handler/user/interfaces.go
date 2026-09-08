package user

import (
	"context"

	"awesome-chat/internal/application/user/command"
	"awesome-chat/internal/application/user/query"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/handler_mocks.go -package=mocks

type useCase interface {
	Register(ctx context.Context, in command.RegisterIn) (command.RegisterOut, error)
	Login(ctx context.Context, in command.LoginIn) (command.LoginOut, error)
	Read(ctx context.Context, in query.ReadIn) (query.ReadOut, error)
	List(ctx context.Context, in query.ListIn) (query.ListOut, error)
	ChatIDs(ctx context.Context, in query.ChatIDsIn) (query.ChatIDsOut, error)
}
