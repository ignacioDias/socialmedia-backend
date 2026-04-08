package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

type LikeRepository interface {
	CreateLike(ctx context.Context, like *models.Like) error
	GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error)
	DeleteLikeByID(ctx context.Context, likeID int64) error
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
	query := `INSERT INTO likes (target_id, target_type, user_id) VALUES ($1, $2, $3) RETURNING like_id`
	return lr.db.QueryRowContext(ctx, query, like.TargetID, like.TargetType, like.UserID).Scan(&like.LikeID)
}

func (lr *likeRepository) GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error) {
	query := `SELECT COUNT(*) FROM likes WHERE target_id = $1 AND target_type = $2`
	var cantLikes int64
	err := lr.db.GetContext(ctx, &cantLikes, query, targetID, targetType)
	return cantLikes, err
}

func (lr *likeRepository) DeleteLikeByID(ctx context.Context, likeID int64) error {
	query := `DELETE FROM likes WHERE like_id = $1`
	result, err := lr.db.ExecContext(ctx, query, likeID)
	return CheckErrResult(result, err, ErrLikeNotFound)
}
