package create

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/application/chat/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	"awesome-chat/internal/domain/core/chat/ports"
	"awesome-chat/internal/domain/core/chat/vo"
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
	userPorts "awesome-chat/internal/domain/core/user/ports"
	userVO "awesome-chat/internal/domain/core/user/vo"
)

type ChatCreateUseCase struct {
	log             appPorts.Logger
	txManager       sharedPorts.TransactionManager
	chatCreateStore ports.CreateWithMembersStore
	userValidator   userPorts.UserValidatorStore
}

func NewChatCreateUseCase(
	log appPorts.Logger,
	txManager sharedPorts.TransactionManager,
	store ports.CreateWithMembersStore,
	userValidator userPorts.UserValidatorStore,
) *ChatCreateUseCase {
	return &ChatCreateUseCase{
		log:             log,
		txManager:       txManager,
		chatCreateStore: store,
		userValidator:   userValidator,
	}
}

func (uc *ChatCreateUseCase) Execute(
	ctx context.Context,
	req dto.CreateChatRequest,
) (
	*dto.ChatResponse,
	error,
) {
	const op = "ChatCreateUseCase.Execute"
	chatID := uuid.New()
	withFields := func(args ...any) []any {
		return append([]any{
			"op", op,
			"chat_id", chatID.String(),
			"chat_name", req.Name,
			"member_count", len(req.MemberIDs),
		}, args...)
	}
	uc.log.Info("Attempting to create new chat", withFields()...)

	if err := uc.basicValidation(
		req.Name,
		req.MemberIDs,
		op,
		withFields,
	); err != nil {
		return nil, err
	}

	members, err := uc.prepareMembers(req.MemberIDs, op, withFields)
	if err != nil {
		return nil, err
	}

	ctxWithTx, err := uc.txManager.BeginAndInjectTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction failed: %w", op, err)
	}

	var txErr error
	defer func() {
		if txErr != nil {
			uc.log.Error("Transaction failed, rolling back",
				withFields("error", txErr.Error())...)
			_ = uc.txManager.RollbackTx(ctxWithTx)
		} else {
			uc.log.Debug("Transaction committed successfully", withFields()...)
		}
	}()

	if err = uc.userValidator.ValidateMultiple(ctx, members); err != nil {
		uc.log.Error("Failed to validate users",
			withFields("error", err.Error())...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if txErr = uc.chatCreateStore.CreateChat(ctxWithTx, vo.ChatID(chatID), req.Name); txErr != nil {
		uc.log.Error("Failed to create new chat",
			withFields("error", txErr.Error())...)
		return nil, fmt.Errorf("%s: %w", op, txErr)
	}

	if txErr = uc.chatCreateStore.AddMembers(ctxWithTx, vo.ChatID(chatID), members); txErr != nil {
		uc.log.Error("Failed to add members to chat",
			withFields("error", txErr.Error())...)
		return nil, fmt.Errorf("%s: %w", op, txErr)
	}

	if txErr = uc.txManager.CommitTx(ctxWithTx); txErr != nil {
		uc.log.Error("Failed to commit transaction",
			withFields("error", txErr.Error())...)
		return nil, fmt.Errorf("commit failed: %w", txErr)
	}

	uc.log.Info("Successfully created new chat", withFields()...)

	return &dto.ChatResponse{
		ID:   chatID.String(),
		Name: req.Name,
	}, nil
}

func (uc *ChatCreateUseCase) basicValidation(
	chatName string,
	memberIDs []string,
	op string,
	logFields func(...any) []any,
) error {
	if len(chatName) < 1 {
		uc.log.Error("Failed to create new chat",
			logFields("error", chatErrors.ErrChatShortName.Error())...)
		return fmt.Errorf("%s: %w", op, chatErrors.ErrChatShortName)
	}
	if len(memberIDs) < 1 {
		uc.log.Error("Failed to create new chat",
			logFields("error", chatErrors.ErrChatInvalidMembersLen.Error())...)
		return fmt.Errorf("%s: %w", op, chatErrors.ErrChatInvalidMembersLen)
	}

	return nil
}

func (uc *ChatCreateUseCase) prepareMembers(
	memberIDs []string,
	op string,
	logFields func(...any) []any,
) (userVO.UserIDs, error) {
	unique := make(map[uuid.UUID]struct{})
	for _, idStr := range memberIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			uc.log.Error("Invalid member ID", logFields("member_id", idStr)...)
			return nil, fmt.Errorf("%s: invalid member ID format: %w", op, err)
		}
		if _, exists := unique[id]; exists {
			uc.log.Warn("Duplicate member ID", logFields("member_id", idStr)...)
			continue // or return nil, errors.New("duplicate member ID")
		}
		unique[id] = struct{}{}
	}

	result := make(userVO.UserIDs, 0, len(unique))
	for id := range unique {
		result = append(result, id)
	}
	return result, nil
}
