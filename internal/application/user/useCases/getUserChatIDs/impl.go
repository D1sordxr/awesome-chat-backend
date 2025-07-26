package getUserChatIDs

import (
	"awesome-chat/internal/application/user/dto"
	"awesome-chat/internal/domain/core/user/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserGetChatIDsUseCase struct {
	store ports.GetUserChatIDsStore
}

func NewUserGetChatIDsUseCase(store ports.GetUserChatIDsStore) *UserGetChatIDsUseCase {
	return &UserGetChatIDsUseCase{store: store}
}

func (uc *UserGetChatIDsUseCase) Execute(
	ctx context.Context,
	userID dto.UserID,
) (
	dto.ChatIDs,
	error,
) {
	const op = "UserGetChatIDsUseCase.Execute"

	id, err := uuid.Parse(string(userID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatIDs, err := uc.store.Execute(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chatIDs, nil
}
