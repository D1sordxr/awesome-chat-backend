package command

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/chat/port"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
)

type Create struct {
	repo  port.Repository
	users port.UserValidator
	tx    port.TxManager
}

func NewCreate(repo port.Repository, users port.UserValidator, tx port.TxManager) *Create {
	return &Create{repo: repo, users: users, tx: tx}
}

type CreateIn struct {
	OwnerID   uuid.UUID
	Name      string
	MemberIDs []string
}

type CreateOut struct {
	ChatID uuid.UUID
	Name   string
}

func (c *Create) Handle(ctx context.Context, in CreateIn) (CreateOut, error) {
	if in.Name == "" {
		return CreateOut{}, chatErrors.ErrChatShortName
	}

	members, err := c.members(in)
	if err != nil {
		return CreateOut{}, err
	}

	if err = c.users.ValidateMultiple(ctx, members); err != nil {
		return CreateOut{}, err
	}

	chatID := uuid.New()

	if err = c.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if createErr := c.repo.Create(txCtx, chatID, in.Name); createErr != nil {
			return createErr
		}

		return c.repo.AddMembers(txCtx, chatID, members)
	}); err != nil {
		return CreateOut{}, err
	}

	return CreateOut{ChatID: chatID, Name: in.Name}, nil
}

func (c *Create) members(in CreateIn) ([]uuid.UUID, error) {
	unique := map[uuid.UUID]struct{}{in.OwnerID: {}}
	members := []uuid.UUID{in.OwnerID}

	for _, raw := range in.MemberIDs {
		memberID, err := uuid.Parse(raw)
		if err != nil {
			return nil, chatErrors.ErrChatInvalidMembersLen
		}

		if _, exists := unique[memberID]; exists {
			continue
		}

		unique[memberID] = struct{}{}
		members = append(members, memberID)
	}

	if len(members) < 2 {
		return nil, chatErrors.ErrChatInvalidMembersLen
	}

	return members, nil
}
