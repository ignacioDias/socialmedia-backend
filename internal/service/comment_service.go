package service

import (
	"context"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error)
	DeleteCommentByID(ctx context.Context, commentID, userID int64) error
	GetCommentsFromTarget(ctx context.Context, targetID int64, targetType models.TargetType, limit, offset int) ([]models.Comment, error)
	GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error)
}

type implCommentService struct {
	commentRepo database.CommentRepository
	cache       *cache.Cache
}

type CommentReq struct {
	Content   string `json:"content"`
	ImagePath string `json:"imagePath"`
}

func NewCommentService(commentRepo database.CommentRepository, cache *cache.Cache) CommentService {
	return &implCommentService{
		commentRepo: commentRepo,
		cache:       cache,
	}
}

func (s *implCommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	return s.commentRepo.CreateComment(ctx, comment)
}

func (s *implCommentService) GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error) {
	return s.commentRepo.GetCommentByID(ctx, commentID)
}

func (s *implCommentService) GetCommentsFromTarget(ctx context.Context, targetID int64, targetType models.TargetType, limit, offset int) ([]models.Comment, error) {
	return s.commentRepo.GetCommentsFromTarget(ctx, targetType, targetID, limit, offset)
}

func (s *implCommentService) DeleteCommentByID(ctx context.Context, commentID, userID int64) error {
	return s.commentRepo.DeleteCommentByID(ctx, commentID, userID)
}

func (s *implCommentService) GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error) {
	return s.commentRepo.GetCommentsFromUser(ctx, userID, limit, offset)
}
