package database

import (
	"context"
	"database/sql"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrCommentNotFound = errors.New("Comment not found")

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error)
	GetCommentsFromTarget(ctx context.Context, targetType models.TargetType, targetID int64, limit, offset int) ([]models.Comment, error)
	GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error)
	DeleteCommentByID(ctx context.Context, commentID int) error
}
type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (cr *commentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (target_id, target_type, user_id, content, image_path) VALUES ($1, $2, $3, $4, $5) RETURNING comment_id`
	return cr.db.QueryRowContext(ctx, query, comment.TargetID, comment.TargetType, comment.UserID, comment.Content, comment.ImagePath).Scan(&comment.CommentID)
}

func (cr *commentRepository) GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error) {
	query := `SELECT comment_id, target_id, target_type, user_id, content, created_at FROM comments WHERE comment_id = $1`
	var comment models.Comment
	if err := cr.db.GetContext(ctx, &comment, query, commentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func (cr *commentRepository) GetCommentsFromTarget(ctx context.Context, targetType models.TargetType, targetID int64, limit, offset int) ([]models.Comment, error) {
	query := `SELECT comment_id, target_id, target_type, user_id, content, created_at, image_path FROM comments WHERE target_type = $1 AND target_id = $2 ORDER BY created_at ASC LIMIT $3 OFFSET $4`
	var comments []models.Comment
	if err := cr.db.SelectContext(ctx, &comments, query, targetType, targetID, limit, offset); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error) {
	query := `SELECT comment_id, target_id, target_type, user_id, content, created_at, image_path FROM comments WHERE user_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	var comments []models.Comment
	if err := cr.db.SelectContext(ctx, &comments, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) DeleteCommentByID(ctx context.Context, commentID int) error {
	query := `DELETE FROM comments WHERE comment_id = $1 OR (target_id = $1 AND target_type = $2)`
	result, err := cr.db.ExecContext(ctx, query, commentID, models.CommentTarget)
	return CheckErrResult(result, err, ErrCommentNotFound)
}
