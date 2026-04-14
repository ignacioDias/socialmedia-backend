package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

type LikeRepository interface {
	CreateLike(ctx context.Context, like *models.Like) error
	GetPostsFromUsersLikes(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error)
	DeleteLike(ctx context.Context, like *models.Like) error
}

var ErrLikeNotFound = errors.New("Like not found")

type likeRepository struct {
	db *sqlx.DB
}

func NewLikeRepository(db *sqlx.DB) LikeRepository {
	return &likeRepository{
		db: db,
	}
}

func (lr *likeRepository) CreateLike(ctx context.Context, like *models.Like) error {
	query := `INSERT INTO likes (target_id, target_type, user_id) VALUES ($1, $2, $3)`
	_, err := lr.db.ExecContext(ctx, query, like.TargetID, like.TargetType, like.UserID)
	return err
}

func (lr *likeRepository) GetPostsFromUsersLikes(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	query := `SELECT p.post_id, p.user_id, p.title, p.content, p.image_path, p.created_at FROM likes l INNER JOIN posts p ON l.target_id = p.post_id WHERE l.user_id = $1 AND l.target_type = $2 ORDER BY l.created_at DESC LIMIT $3 OFFSET $4`
	var posts []models.Post
	if err := lr.db.SelectContext(ctx, &posts, query, userID, models.PostTarget, limit, offset); err != nil {
		return nil, err
	}
	return posts, nil
}

func (lr *likeRepository) GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error) {
	query := `SELECT COUNT(*) FROM likes WHERE target_id = $1 AND target_type = $2`
	var cantLikes int64
	err := lr.db.GetContext(ctx, &cantLikes, query, targetID, targetType)
	return cantLikes, err
}

func (lr *likeRepository) DeleteLike(ctx context.Context, like *models.Like) error {
	query := `DELETE FROM likes WHERE target_type = $1 AND target_id = $2 AND user_id = $3`
	result, err := lr.db.ExecContext(ctx, query, like.TargetType, like.TargetID, like.UserID)
	return CheckErrResult(result, err, ErrLikeNotFound)
}
