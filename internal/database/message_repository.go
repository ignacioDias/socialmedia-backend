package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrMessageNotFound = errors.New("message not found")
var ErrConversationNotFound = errors.New("conversation not found")

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	DeleteMessageByID(ctx context.Context, messageID, userID int64) error
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationsFromUserID(ctx context.Context, userID int64) ([]models.Conversation, error)
	GetMessagesFromConversation(ctx context.Context, conversationID int64, limit, offset int) ([]models.Message, error)
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
	tx, err := mr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	query := `INSERT INTO messages (conversation_id, sender_id, image_path, content) VALUES ($1, $2, $3, $4) RETURNING message_id`
	if err := tx.QueryRowContext(ctx, query, message.ConversationID, message.SenderID, message.ImagePath, message.Content).Scan(&message.MessageID); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE conversations SET last_message_at = NOW() WHERE conversation_id = $1`, message.ConversationID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (mr *messageRepository) DeleteMessageByID(ctx context.Context, messageID, userID int64) error {
	query := `DELETE FROM messages WHERE message_id = $1 AND sender_id = $2`
	result, err := mr.db.ExecContext(ctx, query, messageID, userID)
	return CheckErrResult(result, err, ErrMessageNotFound)
}

func (mr *messageRepository) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	query := `INSERT INTO conversations (user1_id, user2_id) VALUES ($1, $2) RETURNING conversation_id`
	return mr.db.QueryRowContext(ctx, query, conversation.User1ID, conversation.User2ID).Scan(&conversation.ConversationID)
}

func (mr *messageRepository) GetConversationsFromUserID(ctx context.Context, userID int64) ([]models.Conversation, error) {
	query := `SELECT * FROM conversations WHERE (user1_id = $1 OR user2_id = $1) ORDER BY last_message_at DESC`
	var conversations []models.Conversation
	if err := mr.db.SelectContext(ctx, &conversations, query, userID); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (mr *messageRepository) GetMessagesFromConversation(ctx context.Context, conversationID int64, limit, offset int) ([]models.Message, error) {
	query := `SELECT * FROM messages WHERE conversation_id = $1 LIMIT $2 OFFSET $3 ORDER BY created_at ASC`
	var messages []models.Message
	if err := mr.db.SelectContext(ctx, &messages, query, conversationID, limit, offset); err != nil {
		return nil, err
	}
	return messages, nil
}
