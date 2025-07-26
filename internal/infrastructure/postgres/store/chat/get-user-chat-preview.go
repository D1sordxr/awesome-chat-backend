package chat

import (
	"awesome-chat/internal/domain/core/chat/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type GetUserChatPreviewStore struct {
	executor ports.ExecutorManager
}

func NewGetUserChatPreviewStore(executor ports.ExecutorManager) *GetUserChatPreviewStore {
	return &GetUserChatPreviewStore{executor: executor}
}

func (s *GetUserChatPreviewStore) SetupChatPreviews(
	ctx context.Context,
	userID uuid.UUID,
) (
	previews []entity.ChatPreview,
	err error,
) {
	const op = "chat.GetUserChatPreviewStore.SetupChatPreviews"

	conn := s.executor.GetExecutor(ctx)
	query := `
    SELECT 
        c.id AS chat_id,
        c.chat_name AS chat_name,
        m.content AS last_message_content,
        m.user_id AS last_message_sender_id,
        m.created_at AS last_message_time,
        -- TODO read field and count(read)
        0 AS unread_count
    FROM chats c
    JOIN user_chats uc ON c.id = uc.chat_id
    LEFT JOIN LATERAL (
        SELECT content, user_id, created_at
        FROM messages 
        WHERE chat_id = c.id 
        ORDER BY created_at DESC 
        LIMIT 1
    ) m ON true
    WHERE uc.user_id = $1
`
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cp          entity.ChatPreview
			msgText     pgtype.Text
			msgSenderID pgtype.UUID
			msgTime     pgtype.Timestamp
		)

		if err = rows.Scan(
			&cp.ChatID,
			&cp.Name,
			&msgText,
			&msgSenderID,
			&msgTime,
			&cp.UnreadCount,
		); err != nil {
			return nil, err
		}

		if msgText.Valid {
			cp.LastMessage = entity.MessagePreview{
				Text: msgText.String,
				SenderID: func() uuid.UUID {
					if msgSenderID.Valid {
						return msgSenderID.Bytes
					}
					return uuid.Nil
				}(),
				Timestamp: func() time.Time {
					if msgTime.Valid {
						return msgTime.Time
					}
					return time.Time{}
				}(),
			}
		}

		previews = append(previews, cp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return previews, nil
}

func (s *GetUserChatPreviewStore) ParticipantsForChat(
	ctx context.Context,
	chatID uuid.UUID,
) (
	ps []entity.Participant,
	err error,
) {
	const op = "chat.GetUserChatPreviewStore.ParticipantsForChat"
	query := `
        SELECT 
            u.id,
            u.username
        FROM users u
        JOIN user_chats uc ON u.id = uc.user_id
        WHERE uc.chat_id = $1
    `

	rows, err := s.executor.GetExecutor(ctx).Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p entity.Participant
		if err = rows.Scan(
			&p.UserID,
			&p.Username,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ps = append(ps, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ps, nil
}

func (s *GetUserChatPreviewStore) BatchParticipantsForChats(
	ctx context.Context,
	chatIDs []uuid.UUID,
) (map[uuid.UUID][]entity.Participant, error) {
	const op = "chat.GetUserChatPreviewStore.BatchParticipantsForChats"

	query := `
        SELECT 
            uc.chat_id,
            u.id,
            u.username
        FROM users u
        JOIN user_chats uc ON u.id = uc.user_id
        WHERE uc.chat_id = ANY($1)
        ORDER BY uc.chat_id, u.username
    `

	rows, err := s.executor.GetExecutor(ctx).Query(ctx, query, chatIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]entity.Participant)
	for rows.Next() {
		var (
			chatID      uuid.UUID
			participant entity.Participant
		)
		if err = rows.Scan(
			&chatID,
			&participant.UserID,
			&participant.Username,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result[chatID] = append(result[chatID], participant)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
