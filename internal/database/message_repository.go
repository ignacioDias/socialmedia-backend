package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrMessageNotFound = errors.New("message not found")
var ErrChatNotFound = errors.New("chat not found")

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	DeleteMessageByID(ctx context.Context, messageID, userID int64) error
	CreateChat(ctx context.Context, chat *models.Chat) error
	GetChatsFromUserID(ctx context.Context, userID int64) ([]models.Chat, error)
	GetMessagesFromChat(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error)
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
	query := `
        WITH inserted AS (
            INSERT INTO messages (chat_id, sender_id, image_path, content)
            SELECT $1, $2, $3, $4
            FROM chats
            WHERE chat_id = $1 AND (user1_id = $2 OR user2_id = $2)
            RETURNING message_id, chat_id, sender_id, image_path, content, created_at
        ),
        updated AS (
            UPDATE chats SET last_message_at = NOW()
            WHERE chat_id = (SELECT chat_id FROM inserted)
        )
        SELECT * FROM inserted`

	return mr.db.QueryRowContext(ctx, query,
		message.ChatID, message.SenderID, message.ImagePath, message.Content,
	).Scan(&message.MessageID, &message.ChatID, &message.SenderID, &message.ImagePath, &message.Content, &message.CreatedAt)
}
func (mr *messageRepository) DeleteMessageByID(ctx context.Context, messageID, userID int64) error {
	query := `DELETE FROM messages WHERE message_id = $1 AND sender_id = $2`
	result, err := mr.db.ExecContext(ctx, query, messageID, userID)
	return CheckErrResult(result, err, ErrMessageNotFound)
}

func (mr *messageRepository) CreateChat(ctx context.Context, chat *models.Chat) error {
	query := `INSERT INTO chats (user1_id, user2_id) VALUES ($1, $2) RETURNING chat_id`
	return mr.db.QueryRowContext(ctx, query, chat.User1ID, chat.User2ID).Scan(&chat.ChatID)
}

func (mr *messageRepository) GetChatsFromUserID(ctx context.Context, userID int64) ([]models.Chat, error) {
	query := `SELECT * FROM chats WHERE (user1_id = $1 OR user2_id = $1) ORDER BY last_message_at DESC`
	var chats []models.Chat
	if err := mr.db.SelectContext(ctx, &chats, query, userID); err != nil {
		return nil, err
	}
	return chats, nil
}

func (mr *messageRepository) GetMessagesFromChat(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error) {
	query := `SELECT m.* FROM messages m 
    	INNER JOIN chats c ON c.chat_id = m.chat_id 
    	WHERE c.chat_id = $1 AND (c.user1_id = $2 OR c.user2_id = $2) 
    	ORDER BY m.created_at ASC
    	LIMIT $3 OFFSET $4`
	var messages []models.Message
	if err := mr.db.SelectContext(ctx, &messages, query, chatID, userID, limit, offset); err != nil {
		return nil, err
	}
	return messages, nil
}
