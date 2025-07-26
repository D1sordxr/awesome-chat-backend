package authJWT

import (
	"awesome-chat/internal/application/user/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserAuthJWTUseCase struct {
	log         appPorts.Logger
	tokenParser ports.TokenParser
	provider    ports.UserProviderStore
}

func NewUserAuthJWTUseCase(
	log appPorts.Logger,
	tokenParser ports.TokenParser,
	provider ports.UserProviderStore,
) *UserAuthJWTUseCase {
	return &UserAuthJWTUseCase{
		log:         log,
		tokenParser: tokenParser,
		provider:    provider,
	}
}

func (uc *UserAuthJWTUseCase) Execute(
	ctx context.Context,
	tokenStr vo.JWTToken,
) (
	dto.User,
	error,
) {
	const op = "UserAuthJWTUseCase.SetupChatPreviews"
	withFields := func(args ...any) []any { return []any{"op", op, "token", string(tokenStr)} }

	uc.log.Info("Attempting to auth user", withFields()...)

	id, email, err := uc.tokenParser.Do(tokenStr)
	if err != nil {
		uc.log.Error("Failed to parse token", withFields("error", err.Error())...)
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = uuid.Validate(string(id)); err != nil {
		uc.log.Error("Failed to validate token id", withFields("error", err.Error())...)
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.provider.Get(ctx, string(email))
	if err != nil {
		uc.log.Error("Failed to read user",
			withFields("email", string(email), "error", err.Error())...)
		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully auth user", withFields()...)
	return dto.User{
		UserID:   user.UserID.String(),
		Username: user.Username,
		Email:    user.Email.String(),
	}, nil
}
