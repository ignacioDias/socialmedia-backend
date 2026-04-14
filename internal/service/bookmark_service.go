package service

import (
	"context"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
)

type BookmarkService interface {
	CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error
	DeleteBookmark(ctx context.Context, bookmark *models.Bookmark) error
	GetPostsFromUsersBookmarks(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
}

type implBookmarkService struct {
	bookmarkRepo database.BookmarkRepository
	cache        *cache.Cache
}

func NewBookmarkService(bookmarkRepo database.BookmarkRepository, cache *cache.Cache) BookmarkService {
	return &implBookmarkService{
		bookmarkRepo: bookmarkRepo,
		cache:        cache,
	}
}

func (s *implBookmarkService) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	return s.bookmarkRepo.CreateBookmark(ctx, bookmark)
}

func (s *implBookmarkService) DeleteBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	return s.bookmarkRepo.DeleteBookmark(ctx, bookmark)
}

func (s *implBookmarkService) GetPostsFromUsersBookmarks(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	return s.bookmarkRepo.GetPostsFromUsersBookmark(ctx, userID, limit, offset)
}
