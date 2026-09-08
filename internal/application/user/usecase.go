package user

//go:generate ifacemaker -f usecase.go -s UseCaseImpl -i UseCase -p user -o usecase_interface.go

//go:generate gowrap gen -g -p . -i UseCase -t $GOWRAP_OP_TPL -o usecase_with_tracing.go -v "OpPrefix=user.UseCase"

import (
	"context"

	"awesome-chat/internal/application/user/command"
	"awesome-chat/internal/application/user/port"
	"awesome-chat/internal/application/user/query"
)

type UseCaseImpl struct {
	register *command.Register
	login    *command.Login

	read    *query.Read
	list    *query.List
	chatIDs *query.ChatIDs
}

func NewUseCase(
	repo port.Repository,
	provider port.Provider,
	reader port.Reader,
	chatIDsReader port.ChatIDsReader,
	tokens port.TokenCreator,
) *UseCaseImpl {
	return &UseCaseImpl{
		register: command.NewRegister(repo),
		login:    command.NewLogin(provider, tokens),

		read:    query.NewRead(reader),
		list:    query.NewList(reader),
		chatIDs: query.NewChatIDs(chatIDsReader),
	}
}

func (u *UseCaseImpl) Register(ctx context.Context, in command.RegisterIn) (command.RegisterOut, error) {
	return u.register.Handle(ctx, in)
}

func (u *UseCaseImpl) Login(ctx context.Context, in command.LoginIn) (command.LoginOut, error) {
	return u.login.Handle(ctx, in)
}

func (u *UseCaseImpl) Read(ctx context.Context, in query.ReadIn) (query.ReadOut, error) {
	return u.read.Handle(ctx, in)
}

func (u *UseCaseImpl) List(ctx context.Context, in query.ListIn) (query.ListOut, error) {
	return u.list.Handle(ctx, in)
}

func (u *UseCaseImpl) ChatIDs(ctx context.Context, in query.ChatIDsIn) (query.ChatIDsOut, error) {
	return u.chatIDs.Handle(ctx, in)
}
