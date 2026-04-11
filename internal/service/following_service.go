package service

import (
	"context"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type FollowingService interface {
	CreateFollow(ctx context.Context, follow *models.Follow) error
	DeleteFollow(ctx context.Context, follow *models.Follow) error
	GetFollowingFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error)
	GetFollowersFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error)
	GetCountFollowing(ctx context.Context, userID int64) (int, error)
	GetCountFollowers(ctx context.Context, userID int64) (int, error)
}

type implFollowingService struct {
	followingRepo database.FollowRepository
	cache         *cache.Cache
}

func NewFollowingService(followingRepo database.FollowRepository, cache *cache.Cache) FollowingService {
	return &implFollowingService{
		followingRepo: followingRepo,
		cache:         cache,
	}
}

func (s *implFollowingService) CreateFollow(ctx context.Context, follow *models.Follow) error {
	return s.followingRepo.CreateFollow(ctx, follow)
}

func (s *implFollowingService) DeleteFollow(ctx context.Context, follow *models.Follow) error {
	return s.followingRepo.DeleteFollow(ctx, follow)
}

func (s *implFollowingService) GetFollowersFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	return s.followingRepo.GetFollowersFromUserID(ctx, userID, limit, offset)
}

func (s *implFollowingService) GetFollowingFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	return s.followingRepo.GetFollowingFromUserID(ctx, userID, limit, offset)
}

func (s *implFollowingService) GetCountFollowers(ctx context.Context, userID int64) (int, error) {
	return s.followingRepo.GetCantFollowersFromUserID(ctx, userID)
}

func (s *implFollowingService) GetCountFollowing(ctx context.Context, userID int64) (int, error) {
	return s.followingRepo.GetCantFollowingsFromUserID(ctx, userID)
}
