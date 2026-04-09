package database

import (
	"context"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
}

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (mr *messageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	query := `INSERT INTO messages (sender_id, receiver_id, image_path, content) VALUES ($1, $2, $3, $4) RETURNING message_id`
	return mr.db.QueryRowContext(ctx, query, message.SenderID, message.ReceiverID, message.ImagePath, message.Content).Scan(&message.MessageID)
}

func (mr *messageRepository) DeleteMessageByID(ctx context.Context, messageID int64) error {
	result, err := mr.db.ExecContext(ctx, query, messageID)
	return CheckErrResult(result, err, ErrPostNotFound)
}
