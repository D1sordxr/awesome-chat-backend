package create

import (
	"awesome-chat/internal/application/chat/dto"
	"awesome-chat/internal/domain/core/chat/ports"
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
	userPorts "awesome-chat/internal/domain/core/user/ports"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ChatCreateUseCase struct {
	txManager      sharedPorts.TransactionManager
	store          ports.CreateWithMembersStore
	userValidation userPorts.UserGetAllStore
}

func NewChatCreateUseCase(
	txManager sharedPorts.TransactionManager,
	store ports.CreateWithMembersStore,
	userValidation userPorts.UserGetAllStore,
) *ChatCreateUseCase {
	return &ChatCreateUseCase{
		txManager:      txManager,
		store:          store,
		userValidation: userValidation,
	}
}

func (uc *ChatCreateUseCase) Execute(ctx context.Context, req dto.CreateChatRequest) (*dto.ChatResponse, error) {
	// basic validation
	if len(req.Name) < 1 {
		return nil, errors.New("name is too short")
	}
	if len(req.MemberIDs) < 1 {
		return nil, errors.New("chat requires at least one member")
	}

	// duplicate check
	memberSet := make(map[uuid.UUID]struct{}, len(req.MemberIDs))
	for _, memberID := range req.MemberIDs {
		id, err := uuid.Parse(memberID)
		if err != nil {
			return nil, errors.New("invalid member ID format")
		}
		if _, exists := memberSet[id]; exists {
			return nil, errors.New("duplicate member ID")
		}
		memberSet[id] = struct{}{}
	}

	ctxWithTx, err := uc.txManager.BeginAndInjectTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}

	var txErr error
	defer func() {
		if txErr != nil {
			_ = uc.txManager.RollbackTx(ctxWithTx)
		}
	}()

	chatID := uuid.New()
	validatedMembers := make([]uuid.UUID, 0, len(memberSet))
	errGroup, ctx := errgroup.WithContext(ctx)

	// user validation goroutine
	errGroup.Go(func() error {
		for id := range memberSet {
			if _, err := uc.userValidation.GetAll(ctx); err != nil { // TODO change.
				return fmt.Errorf("user validation failed for user %s: %w", id, err)
			}
			validatedMembers = append(validatedMembers, id)
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		txErr = err
		return nil, err
	}

	if err := uc.store.Execute(ctxWithTx, chatID, req.Name); err != nil {
		txErr = err
		return nil, err
	}

	for _, memberID := range validatedMembers {
		if err := uc.store.AddMember(ctxWithTx, chatID, memberID); err != nil {
			txErr = err
			return nil, err
		}
	}

	if err := uc.txManager.CommitTx(ctxWithTx); err != nil {
		txErr = err
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return &dto.ChatResponse{
		ID:   chatID.String(),
		Name: req.Name,
	}, nil
}
