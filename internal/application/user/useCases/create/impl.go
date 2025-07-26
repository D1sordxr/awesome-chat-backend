package create

import (
	"awesome-chat/internal/application/user/dto"
	"awesome-chat/internal/domain/core/user/ports"
	"context"
	"errors"

	"github.com/google/uuid"
)

type UserCreateUseCase struct {
	store ports.UserCreateStore
}

func NewUserCreateUseCase(store ports.UserCreateStore) *UserCreateUseCase {
	return &UserCreateUseCase{store: store}
}

func (u *UserCreateUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	if len(req.Username) < 1 {
		return dto.UserResponse{}, errors.New("username is too short")
	}

	id := uuid.New()
	if err := u.store.Execute(ctx, id, req.Username); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{ID: id.String()}, nil
}
