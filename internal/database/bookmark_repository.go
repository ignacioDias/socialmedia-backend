package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrBookmarkNotFound = errors.New("Bookmark not found")

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error
	DeleteBookmark(ctx context.Context, bookmark *models.Bookmark) error
	GetPostsFromUsersBookmark(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
}

type bookmarkRepository struct {
	db *sqlx.DB
}

func NewBookmarkRepository(db *sqlx.DB) BookmarkRepository {
	return &bookmarkRepository{
		db: db,
	}
}

func (br *bookmarkRepository) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	query := `INSERT INTO bookmarks (post_id, user_id) VALUES ($1, $2)`
	_, err := br.db.ExecContext(ctx, query, bookmark.PostID, bookmark.UserID)
	return err
}

func (br *bookmarkRepository) GetPostsFromUsersBookmark(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	query := `SELECT p.post_id, p.user_id, p.title, p.content, p.image_path, p.created_at 
          FROM posts p 
          INNER JOIN bookmarks b ON p.post_id = b.post_id 
          WHERE b.user_id = $1
		  ORDER BY b.created_at DESC
		  LIMIT $2 OFFSET $3
		  `
	var posts []models.Post
	if err := br.db.SelectContext(ctx, &posts, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return posts, nil
}

func (br *bookmarkRepository) DeleteBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	query := `DELETE FROM bookmarks WHERE post_id = $1 AND user_id = $2`
	result, err := br.db.ExecContext(ctx, query, bookmark.PostID, bookmark.UserID)
	return CheckErrResult(result, err, ErrBookmarkNotFound)
}
