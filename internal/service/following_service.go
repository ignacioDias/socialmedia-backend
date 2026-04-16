package service

import (
	"context"
	"fmt"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
	"time"
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

var ErrSelfFollow = fmt.Errorf("cannot follow yourself")

func NewFollowingService(followingRepo database.FollowRepository, cache *cache.Cache) FollowingService {
	return &implFollowingService{
		followingRepo: followingRepo,
		cache:         cache,
	}
}

func (s *implFollowingService) CreateFollow(ctx context.Context, follow *models.Follow) error {
	if follow.FollowerID == follow.FollowingID {
		return ErrSelfFollow
	}
	if err := s.followingRepo.CreateFollow(ctx, follow); err != nil {
		return err
	}
	followersCountKey := fmt.Sprintf("followersCount:%d", follow.FollowingID)
	followingCountKey := fmt.Sprintf("followingCount:%d", follow.FollowerID)
	_ = s.cache.Delete(followersCountKey)
	_ = s.cache.Delete(followingCountKey)
	return nil
}

func (s *implFollowingService) DeleteFollow(ctx context.Context, follow *models.Follow) error {
	if err := s.followingRepo.DeleteFollow(ctx, follow); err != nil {
		return err
	}
	followersCountKey := fmt.Sprintf("followersCount:%d", follow.FollowingID)
	followingCountKey := fmt.Sprintf("followingCount:%d", follow.FollowerID)
	_ = s.cache.Delete(followersCountKey)
	_ = s.cache.Delete(followingCountKey)
	return nil
}

func (s *implFollowingService) GetFollowersFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	return s.followingRepo.GetFollowersFromUserID(ctx, userID, limit, offset)
}

func (s *implFollowingService) GetFollowingFromID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	return s.followingRepo.GetFollowingFromUserID(ctx, userID, limit, offset)
}

func (s *implFollowingService) GetCountFollowers(ctx context.Context, userID int64) (int, error) {
	if userID <= 0 {
		return 0, ErrInvalidID
	}
	key := fmt.Sprintf("followersCount:%d", userID)
	var followersCount int
	if err := s.cache.Get(key, &followersCount); err == nil {
		return followersCount, nil
	}

	followersCount, err := s.followingRepo.GetCantFollowersFromUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	_ = s.cache.Set(key, followersCount, time.Hour)
	return followersCount, nil
}

func (s *implFollowingService) GetCountFollowing(ctx context.Context, userID int64) (int, error) {
	if userID <= 0 {
		return 0, ErrInvalidID
	}
	key := fmt.Sprintf("followingCount:%d", userID)
	var followingCount int
	if err := s.cache.Get(key, &followingCount); err == nil {
		return followingCount, nil
	}

	followingCount, err := s.followingRepo.GetCantFollowingsFromUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	_ = s.cache.Set(key, followingCount, time.Hour)
	return followingCount, nil
}
