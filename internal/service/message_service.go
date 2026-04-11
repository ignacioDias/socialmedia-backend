package service

import (
	"context"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type MessageService interface {
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationsFromUser(ctx context.Context, userID int64) ([]models.Conversation, error)
	DeleteMessage(ctx context.Context, messageID, userID int64) error
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

func (s *implMessageService) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	return s.messageRepo.CreateConversation(ctx, conversation)
}

func (s *implMessageService) GetConversationsFromUser(ctx context.Context, userID int64) ([]models.Conversation, error) {
	return s.messageRepo.GetConversationsFromUserID(ctx, userID)
}

func (s *implMessageService) DeleteMessage(ctx context.Context, messageID, userID int64) error {
	return s.messageRepo.DeleteMessageByID(ctx, messageID, userID)
}
