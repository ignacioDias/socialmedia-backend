package service

import (
	"context"
	"fmt"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
	"time"
)

type PostRequest struct {
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	ImagePath string `db:"image_path" json:"imagePath"`
}

type PostService interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetPostFromID(ctx context.Context, postID int64) (*models.Post, error)
	DeletePost(ctx context.Context, postID, userID int64) error
	GetPostsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	GetPostsFromFollowing(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
}

type implPostService struct {
	postRepo database.PostRepository
	cache    *cache.Cache
}

func NewPostService(postRepo database.PostRepository, cache *cache.Cache) PostService {
	return &implPostService{
		postRepo: postRepo,
		cache:    cache,
	}
}

func (s *implPostService) CreatePost(ctx context.Context, post *models.Post) error {
	return s.postRepo.CreatePost(ctx, post)
}

func (s *implPostService) GetPostFromID(ctx context.Context, postID int64) (*models.Post, error) {
	if postID <= 0 {
		return nil, ErrInvalidID
	}
	key := fmt.Sprintf("post:%d", postID)
	var post *models.Post
	if err := s.cache.Get(key, &post); err == nil {
		return post, nil
	}
	post, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(key, post, time.Hour)
	return post, nil
}
func (s *implPostService) DeletePost(ctx context.Context, postID, userID int64) error {
	if userID <= 0 || postID <= 0 {
		return ErrInvalidID
	}
	if err := s.postRepo.DeletePostByID(ctx, postID, userID); err != nil {
		return err
	}
	key := fmt.Sprintf("post:%d", postID)
	_ = s.cache.Delete(key)
	return nil
}

func (s *implPostService) GetPostsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	return s.postRepo.GetPostsByUserID(ctx, userID, limit, offset)
}

func (s *implPostService) GetPostsFromFollowing(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, ErrInvalidID
	}
	return s.postRepo.GetPostsFromFollows(ctx, userID, limit, offset)
}
