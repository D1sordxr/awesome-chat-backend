package message

import (
	"context"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/application/message/query"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/handler_mocks.go -package=mocks

type useCase interface {
	Send(ctx context.Context, in command.SendIn) (command.SendOut, error)
	CreateVoiceUpload(ctx context.Context, in command.CreateVoiceUploadIn) (command.CreateVoiceUploadOut, error)
	SendVoice(ctx context.Context, in command.SendVoiceIn) (command.SendVoiceOut, error)
	List(ctx context.Context, in query.ListIn) (query.ListOut, error)
}
