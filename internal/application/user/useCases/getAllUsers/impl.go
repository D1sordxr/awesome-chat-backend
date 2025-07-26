package getAllUsers

import (
	"awesome-chat/internal/application/user/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/ports"
	"context"
	"fmt"
)

type UsersGetAllUseCase struct {
	log   appPorts.Logger
	store ports.UserGetAllStore
}

func NewUsersGetAllUseCase(
	log appPorts.Logger,
	store ports.UserGetAllStore,
) *UsersGetAllUseCase {
	return &UsersGetAllUseCase{
		log:   log,
		store: store,
	}
}

func (uc *UsersGetAllUseCase) Execute(
	ctx context.Context,
) (
	dto.GetAllUsersResponse,
	error,
) {
	const op = "UsersGetAllUseCase.Execute"
	withFields := func(args ...any) []any {
		return append([]any{"op", op}, args...)
	}

	uc.log.Info("Attempting to get all users", withFields()...)

	users, err := uc.store.GetAll(ctx)
	if err != nil {
		uc.log.Error("Failed to get all users", withFields("error", err.Error())...)
		return dto.GetAllUsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	usersResponse := make(dto.Users, 0, len(users))
	for _, user := range users {
		usersResponse = append(usersResponse, dto.User{
			UserID:   user.UserID.String(),
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return dto.GetAllUsersResponse{Users: usersResponse}, nil
}
