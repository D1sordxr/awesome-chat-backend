package register

import (
	"awesome-chat/internal/application/user/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterUseCase struct {
	log  appPorts.Logger
	repo ports.UserRepo
}

func NewUserRegisterUseCase(
	log appPorts.Logger,
	repo ports.UserRepo,
) *UserRegisterUseCase {
	return &UserRegisterUseCase{
		log:  log,
		repo: repo,
	}
}

func (uc *UserRegisterUseCase) Execute(
	ctx context.Context,
	req dto.RegisterUserRequest,
) (
	dto.RegisterUserResponse,
	error,
) {
	const op = "UserRegisterUseCase.SetupChatPreviews"
	withFields := func(args ...any) []any { return append([]any{"op", op, "email", req.Email}, args...) }

	uc.log.Info("Attempting to register user", withFields()...)

	email, err := vo.NewEmail(req.Email)
	if err != nil {
		uc.log.Error("Failed to create new email", withFields("error", err.Error())...)
		return dto.RegisterUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Error("Failed to generate password", withFields("error", err.Error())...)
		return dto.RegisterUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	userID := uuid.New()

	userEntity := entity.NewUser(
		userID,
		req.Username,
		email,
		pass,
	)

	if err = uc.repo.Save(ctx, userEntity); err != nil {
		uc.log.Error("Failed to save user", withFields("error", err.Error())...)
		return dto.RegisterUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("User registered successfully", withFields()...)
	return dto.RegisterUserResponse{
		UserID: userID.String(),
	}, nil
}
