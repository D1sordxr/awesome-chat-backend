package command

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/chat/port"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
)

type AddMember struct {
	repo  port.Repository
	chats port.Validator
	users port.UserValidator
}

func NewAddMember(repo port.Repository, chats port.Validator, users port.UserValidator) *AddMember {
	return &AddMember{repo: repo, chats: chats, users: users}
}

type AddMemberIn struct {
	CallerID uuid.UUID
	ChatID   uuid.UUID
	UserID   uuid.UUID
}

type AddMemberOut struct{}

func (c *AddMember) Handle(ctx context.Context, in AddMemberIn) (AddMemberOut, error) {
	callerIsMember, err := c.chats.IsMember(ctx, in.ChatID, in.CallerID)
	if err != nil {
		return AddMemberOut{}, err
	}

	if !callerIsMember {
		return AddMemberOut{}, chatErrors.ErrChatAccessDenied
	}

	if err = c.users.ValidateByID(ctx, in.UserID); err != nil {
		return AddMemberOut{}, err
	}

	alreadyMember, err := c.chats.IsMember(ctx, in.ChatID, in.UserID)
	if err != nil {
		return AddMemberOut{}, err
	}

	if alreadyMember {
		return AddMemberOut{}, chatErrors.ErrChatMemberExists
	}

	if err = c.repo.AddMember(ctx, in.ChatID, in.UserID); err != nil {
		return AddMemberOut{}, err
	}

	return AddMemberOut{}, nil
}
