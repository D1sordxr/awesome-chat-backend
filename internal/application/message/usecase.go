package message

//go:generate ifacemaker -f usecase.go -s UseCaseImpl -i UseCase -p message -o usecase_interface.go

//go:generate gowrap gen -g -p . -i UseCase -t $GOWRAP_OP_TPL -o usecase_with_tracing.go -v "OpPrefix=message.UseCase"

import (
	"context"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/application/message/port"
	"awesome-chat/internal/application/message/query"
)

type UseCaseImpl struct {
	send        *command.Send
	voiceUpload *command.CreateVoiceUpload
	sendVoice   *command.SendVoice

	list *query.List
}

func NewUseCase(
	repo port.Repository,
	voiceRepo port.VoiceRepository,
	outbox port.Outbox,
	reader port.Reader,
	uploads port.UploadService,
	chats port.ChatValidator,
	tx port.TxManager,
) *UseCaseImpl {
	return &UseCaseImpl{
		send:        command.NewSend(repo, outbox, chats, tx),
		voiceUpload: command.NewCreateVoiceUpload(uploads, chats),
		sendVoice:   command.NewSendVoice(voiceRepo, chats, tx),

		list: query.NewList(reader, chats),
	}
}

func (u *UseCaseImpl) Send(ctx context.Context, in command.SendIn) (command.SendOut, error) {
	return u.send.Handle(ctx, in)
}

func (u *UseCaseImpl) CreateVoiceUpload(
	ctx context.Context,
	in command.CreateVoiceUploadIn,
) (command.CreateVoiceUploadOut, error) {
	return u.voiceUpload.Handle(ctx, in)
}

func (u *UseCaseImpl) SendVoice(ctx context.Context, in command.SendVoiceIn) (command.SendVoiceOut, error) {
	return u.sendVoice.Handle(ctx, in)
}

func (u *UseCaseImpl) List(ctx context.Context, in query.ListIn) (query.ListOut, error) {
	return u.list.Handle(ctx, in)
}
