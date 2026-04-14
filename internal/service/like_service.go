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

const likesCountCacheKeyFormat = "likes:count:%s:%d"

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
	key := fmt.Sprintf(likesCountCacheKeyFormat, like.TargetType, like.TargetID)
	if err := s.likeRepo.CreateLike(ctx, like); err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}
	_ = s.cache.Delete(key)
	return nil
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
	key := fmt.Sprintf(likesCountCacheKeyFormat, targetType, targetID)
	var likeCount int64
	if err := s.cache.Get(key, &likeCount); err == nil {
		return likeCount, nil
	}
	likeCount, err := s.likeRepo.GetLikesCountFromTarget(ctx, targetID, targetType)
	if err != nil {
		return 0, fmt.Errorf("failed to get likes count: %w", err)
	}
	_ = s.cache.Set(key, likeCount, 5*time.Minute)
	return likeCount, nil
}

func (s *implLikeService) DeleteLike(ctx context.Context, like *models.Like) error {
	if like == nil {
		return errors.New("invalid like")
	}
	key := fmt.Sprintf(likesCountCacheKeyFormat, like.TargetType, like.TargetID)
	if err := s.likeRepo.DeleteLike(ctx, like); err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}
	_ = s.cache.Delete(key)
	return nil
}
