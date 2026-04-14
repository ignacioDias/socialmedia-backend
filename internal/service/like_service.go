package service

import (
	"context"
	"errors"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type LikeService interface {
	CreateLike(ctx context.Context, like *models.Like) error
	GetPostsFromUsersLikes(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error)
	DeleteLike(ctx context.Context, like *models.Like) error
}

type implLikeService struct {
	likeRepo database.LikeRepository
	cache    *cache.Cache
}

func NewLikeService(likeRepo database.LikeRepository, cache *cache.Cache) LikeService {
	return &implLikeService{
		likeRepo: likeRepo,
		cache:    cache,
	}
}

func (s *implLikeService) CreateLike(ctx context.Context, like *models.Like) error {
	return s.likeRepo.CreateLike(ctx, like)
}

func (s *implLikeService) GetPostsFromUsersLikes(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.likeRepo.GetPostsFromUsersLikes(ctx, userID, limit, offset)
}

func (s *implLikeService) GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error) {
	if targetID <= 0 {
		return 0, errors.New("invalid target ID")
	}
	return s.likeRepo.GetLikesCountFromTarget(ctx, targetID, targetType)
}

func (s *implLikeService) DeleteLike(ctx context.Context, like *models.Like) error {
	if like == nil {
		return errors.New("invalid like")
	}
	return s.likeRepo.DeleteLike(ctx, like)
}
