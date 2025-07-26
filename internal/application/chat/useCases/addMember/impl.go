package addMember

import (
	"awesome-chat/internal/application/chat/dto"
	"awesome-chat/internal/domain/core/chat/ports"
	userPorts "awesome-chat/internal/domain/core/user/ports"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ChatAddMemberUseCase struct {
	chatStore     ports.AddMemberStore
	chatValidator ports.ValidateStore
	userValidator userPorts.UserGetAllStore
}

func NewChatAddMemberUseCase(
	chatStore ports.AddMemberStore,
	chatValidator ports.ValidateStore,
	userValidator userPorts.UserGetAllStore,
) *ChatAddMemberUseCase {
	return &ChatAddMemberUseCase{
		chatStore:     chatStore,
		chatValidator: chatValidator,
		userValidator: userValidator,
	}
}

func (uc *ChatAddMemberUseCase) Execute(ctx context.Context, req dto.AddUserRequest) error {
	chatID, err := uuid.Parse(req.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err = uc.chatValidator.ValidateExists(groupCtx, chatID); err != nil {
			return fmt.Errorf("chat validation failed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if isMember, validateErr := uc.chatValidator.IsMember(groupCtx, chatID, userID); validateErr != nil {
			return fmt.Errorf("chat is member: %w", validateErr)
		} else if isMember {
			return fmt.Errorf("user already in chat") // 409
		}
		return nil
	})

	g.Go(func() error {
		if _, err = uc.userValidator.GetAll(groupCtx); err != nil { // TODO change
			return fmt.Errorf("user validation failed: %w", err)
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		return err
	}

	if err = uc.chatStore.AddMember(ctx, chatID, userID); err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}
