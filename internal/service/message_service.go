package service

import (
	"context"
	"errors"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type MessageService interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	GetChatsFromUser(ctx context.Context, userID int64) ([]models.Chat, error)
	DeleteMessage(ctx context.Context, messageID, userID int64) error
	GetMessagesFromChat(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error)
	CreateMessage(ctx context.Context, message *models.Message) error
}

type MessageRequest struct {
	ImagePath string `json:"imagePath"`
	Content   string `json:"content"`
}

type implMessageService struct {
	messageRepo database.MessageRepository
	cache       *cache.Cache
}

func NewMessageService(messageRepo database.MessageRepository, cache *cache.Cache) MessageService {
	return &implMessageService{
		messageRepo: messageRepo,
		cache:       cache,
	}
}

func (s *implMessageService) CreateChat(ctx context.Context, chat *models.Chat) error {
	if chat.User1ID == chat.User2ID {
		return errors.New("can't make a chat with yourself")
	}
	return s.messageRepo.CreateChat(ctx, chat)
}

func (s *implMessageService) GetChatsFromUser(ctx context.Context, userID int64) ([]models.Chat, error) {
	return s.messageRepo.GetChatsFromUserID(ctx, userID)
}

func (s *implMessageService) DeleteMessage(ctx context.Context, messageID, userID int64) error {
	return s.messageRepo.DeleteMessageByID(ctx, messageID, userID)
}

func (s *implMessageService) GetMessagesFromChat(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error) {
	return s.messageRepo.GetMessagesFromChat(ctx, chatID, userID, limit, offset)
}

func (s *implMessageService) CreateMessage(ctx context.Context, message *models.Message) error {
	if message.Content == "" && message.ImagePath == "" {
		return errors.New("empty message")
	}
	return s.messageRepo.CreateMessage(ctx, message)
}
