package user

import (
	"awesome-chat/internal/application/user/dto"
	userErrors "awesome-chat/internal/domain/core/user/errors"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	registerUseCase interface {
		Execute(ctx context.Context, req dto.RegisterUserRequest) (dto.RegisterUserResponse, error)
	}
	loginUseCase interface {
		Execute(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	}
	authJWTUseCase interface {
		Execute(ctx context.Context, token vo.JWTToken) (dto.User, error)
	}
	getUserChatIDsUseCase interface {
		Execute(ctx context.Context, userID dto.UserID) (dto.ChatIDs, error)
	}
	getAllUsersUseCase interface {
		Execute(ctx context.Context) (dto.GetAllUsersResponse, error)
	}

	// V1 useCases
	//createUseCase interface {
	//	SetupChatPreviews(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	//}
)

type Handler struct {
	registerUC       registerUseCase
	loginUC          loginUseCase
	authJWTUC        authJWTUseCase
	getUserChatIDsUC getUserChatIDsUseCase
	getAllUsersUC    getAllUsersUseCase

	//createUC createUseCase
}

func NewUserHandler(
	registerUC registerUseCase,
	loginUC loginUseCase,
	authJWTUC authJWTUseCase,
	getUserChatIDsUC getUserChatIDsUseCase,
	getAllUsersUC getAllUsersUseCase,
// createUC createUseCase,
// getUC getUseCase,
) *Handler {
	return &Handler{
		registerUC:       registerUC,
		loginUC:          loginUC,
		authJWTUC:        authJWTUC,
		getUserChatIDsUC: getUserChatIDsUC,
		getAllUsersUC:    getAllUsersUC,
		//createUC:   createUC,
		//getUC:      getUC,
	}
}

func (h *Handler) register(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	var req dto.RegisterUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	resp, err := h.registerUC.Execute(reqCtx, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled):
			return ctx.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
				"error": "Request timeout",
			})
		case errors.Is(err, userErrors.ErrInvalidEmailLength), errors.Is(err, userErrors.ErrInvalidEmailFormat):
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":         "Invalid email format",
				"valid_example": "user@example.com",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      resp.UserID,
		"message": "Registration successful",
	})
}

func (h *Handler) login(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	resp, err := h.loginUC.Execute(reqCtx, req)
	if err != nil {
		switch {
		case errors.Is(err, userErrors.ErrUserDoesNotExist):
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    resp.Token,
		Expires:  time.Now().Add(time.Hour * 48),
		HTTPOnly: true,
		// Secure:   true,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": resp.Username,
		"token":    resp.Token,
		"message":  "Login successful",
	})
}

func (h *Handler) logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func (h *Handler) authJWT(ctx *fiber.Ctx) error {
	tokenStr := ctx.Cookies("jwt")
	if tokenStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authorization token required",
			"details": "JWT cookie is missing",
		})
	}

	user, err := h.authJWTUC.Execute(ctx.Context(), vo.JWTToken(tokenStr))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid or expired token",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.UserID,
		"username": user.Username,
		"message":  "Authentication successful",
	})
}

func (h *Handler) getChatIDs(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	userID := ctx.Params("id", "")
	if userID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id required",
		})
	}

	resp, err := h.getUserChatIDsUC.Execute(reqCtx, dto.UserID(userID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get chat ids",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"chat_ids": resp,
	})
}

func (h *Handler) getAllUsers(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.getAllUsersUC.Execute(reqCtx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": resp.Users,
	})
}

//func (h *Handler) CreateUser(ctx *fiber.Ctx) error {
//	var req dto.CreateUserRequest
//	if err := ctx.BodyParser(&req); err != nil {
//		return fiber.NewError(fiber.StatusBadRequest, err.Error())
//	}
//
//	user, err := h.createUC.SetupChatPreviews(ctx.Context(), req)
//	if err != nil {
//		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
//	}
//
//	return ctx.Status(fiber.StatusCreated).JSON(user)
//}
//
//func (h *Handler) GetUser(ctx *fiber.Ctx) error {
//	userID := ctx.Params("id")
//	if userID == "" {
//		return fiber.NewError(fiber.StatusBadRequest, "user id is required")
//	}
//
//	user, err := h.getUC.SetupChatPreviews(ctx.Context(), userID)
//	if err != nil {
//		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
//	}
//
//	return ctx.JSON(user)
//}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	// V2
	router.Post("/user/register", h.register)
	router.Post("/user/login", h.login)
	router.Get("/user/logout", h.logout)
	router.Get("/user/auth-jwt", h.authJWT)
	router.Get("/user/chat-ids/:id", h.getChatIDs)
	router.Get("/user/get-all", h.getAllUsers)

	// V1
	//router.Post("/user", h.CreateUser)
	//router.Get("/user/:id", h.GetUser)
}
