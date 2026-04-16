package service

import (
	"context"
	"errors"
	"fmt"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
	"time"
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

var ErrInvalidID = errors.New("invalid ID")

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
	if commentID <= 0 {
		return nil, ErrInvalidID
	}
	key := fmt.Sprintf("comment:%d", commentID)
	var comment *models.Comment
	if err := s.cache.Get(key, &comment); err == nil {
		return comment, nil
	}
	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(key, comment, time.Hour)
	return comment, nil
}

func (s *implCommentService) GetCommentsFromTarget(ctx context.Context, targetID int64, targetType models.TargetType, limit, offset int) ([]models.Comment, error) {
	if targetID <= 0 {
		return nil, ErrInvalidID
	}
	return s.commentRepo.GetCommentsFromTarget(ctx, targetType, targetID, limit, offset)
}

func (s *implCommentService) DeleteCommentByID(ctx context.Context, commentID, userID int64) error {
	if userID <= 0 || commentID <= 0 {
		return ErrInvalidID
	}

	if err := s.commentRepo.DeleteCommentByID(ctx, commentID, userID); err != nil {
		return err
	}
	key := fmt.Sprintf("comment:%d", commentID)
	_ = s.cache.Delete(key)
	return nil
}

func (s *implCommentService) GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	return s.commentRepo.GetCommentsFromUser(ctx, userID, limit, offset)
}
