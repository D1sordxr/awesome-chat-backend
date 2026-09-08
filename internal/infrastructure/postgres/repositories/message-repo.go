package repositories

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type MessageRepo struct {
	e ports.ExecutorManager
}

func NewMessageRepo(e ports.ExecutorManager) *MessageRepo {
	return &MessageRepo{e: e}
}

func (r *MessageRepo) SaveBatch(ctx context.Context, messages []entity.Message) error {
	const op = "repositories.MessageRepo.SaveBatch"

	if len(messages) == 0 {
		return nil
	}

	userIDs := make([]uuid.UUID, 0, len(messages))
	chatIDs := make([]uuid.UUID, 0, len(messages))
	contents := make([]string, 0, len(messages))
	timestamps := make([]time.Time, 0, len(messages))

	for _, message := range messages {
		userIDs = append(userIDs, message.UserID)
		chatIDs = append(chatIDs, message.ChatID)
		contents = append(contents, message.Content)
		timestamps = append(timestamps, message.Timestamp)
	}

	query := `
        INSERT INTO messages (
            user_id,
            chat_id,
            content,
            timestamp
        ) SELECT 
            unnest($1::uuid[]),
            unnest($2::uuid[]),
            unnest($3::text[]),
            unnest($4::timestamp[])`

	_, err := r.e.GetExecutor(ctx).Exec(ctx, query, userIDs, chatIDs, contents, timestamps)
	return fmt.Errorf("%s: %w", op, err)
}

func (r *MessageRepo) SaveBatchFast(ctx context.Context, messages []entity.Message) error {
	const op = "repositories.MessageRepo.SaveBatchFast"

	if len(messages) == 0 {
		return nil
	}

	tx, err := r.e.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction failed: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"messages"},
		[]string{"user_id", "chat_id", "content", "timestamp"},
		pgx.CopyFromSlice(len(messages), func(i int) ([]interface{}, error) {
			msg := messages[i]
			return []interface{}{msg.UserID, msg.ChatID, msg.Content, msg.Timestamp}, nil
		}),
	)

	if err != nil {
		return fmt.Errorf("%s: copy from failed: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (r *MessageRepo) Save(ctx context.Context, message entity.Message) (entity.Message, error) {
	const op = "repositories.MessageRepo.Save"

	query := `
		INSERT INTO messages (
			user_id,
			chat_id,
			content
		) VALUES ($1, $2, $3)
		RETURNING id, created_at`

	saved := message

	if err := r.e.GetExecutor(ctx).QueryRow(ctx, query,
		message.UserID,
		message.ChatID,
		message.Content,
	).Scan(&saved.ID, &saved.Timestamp); err != nil {
		return entity.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}
