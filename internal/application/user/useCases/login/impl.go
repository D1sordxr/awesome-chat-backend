package login

import (
	"awesome-chat/internal/application/user/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginUseCase struct {
	log          appPorts.Logger
	provider     ports.UserProviderStore
	tokenCreator ports.TokenCreator
}

func NewUserLoginUseCase(
	log appPorts.Logger,
	provider ports.UserProviderStore,
	creator ports.TokenCreator,
) *UserLoginUseCase {
	return &UserLoginUseCase{
		log:          log,
		provider:     provider,
		tokenCreator: creator,
	}
}

func (uc *UserLoginUseCase) Execute(
	ctx context.Context,
	req dto.LoginRequest,
) (
	resp dto.LoginResponse,
	err error,
) {
	const op = "UserLoginUseCase.SetupChatPreviews"
	withFields := func(args ...any) []any { return append([]any{"op", op, "email", req.Email}, args...) }

	uc.log.Info("Attempting to login user", withFields()...)
	defer func() {
		if err != nil {
			uc.log.Error("Failed to login user", withFields("error", err.Error())...)
		}
	}()

	email, err := vo.NewEmail(req.Email)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.provider.Get(ctx, email.String())
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}

	token, err := uc.tokenCreator.Do(user)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully logged in", withFields()...)
	resp.Token = token
	resp.Username = user.Username
	return resp, nil
}
