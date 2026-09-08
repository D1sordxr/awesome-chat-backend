package chat

//go:generate ifacemaker -f usecase.go -s UseCaseImpl -i UseCase -p chat -o usecase_interface.go

//go:generate gowrap gen -g -p . -i UseCase -t $GOWRAP_OP_TPL -o usecase_with_tracing.go -v "OpPrefix=chat.UseCase"

import (
	"context"

	"awesome-chat/internal/application/chat/command"
	"awesome-chat/internal/application/chat/port"
	"awesome-chat/internal/application/chat/query"
)

type UseCaseImpl struct {
	create    *command.Create
	addMember *command.AddMember

	previews *query.Previews
}

func NewUseCase(
	repo port.Repository,
	chats port.Validator,
	users port.UserValidator,
	reader port.PreviewReader,
	tx port.TxManager,
) *UseCaseImpl {
	return &UseCaseImpl{
		create:    command.NewCreate(repo, users, tx),
		addMember: command.NewAddMember(repo, chats, users),

		previews: query.NewPreviews(reader),
	}
}

func (u *UseCaseImpl) Create(ctx context.Context, in command.CreateIn) (command.CreateOut, error) {
	return u.create.Handle(ctx, in)
}

func (u *UseCaseImpl) AddMember(ctx context.Context, in command.AddMemberIn) (command.AddMemberOut, error) {
	return u.addMember.Handle(ctx, in)
}

func (u *UseCaseImpl) Previews(ctx context.Context, in query.PreviewsIn) (query.PreviewsOut, error) {
	return u.previews.Handle(ctx, in)
}
