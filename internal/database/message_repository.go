package database

import (
	"context"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
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
	return nil
}
