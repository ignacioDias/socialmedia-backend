package database

import (
	"context"
	"database/sql"
	"errors"
	models "socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrPostNotFound = errors.New("Post not found")

type PostRepository interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)
	GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	DeletePostByID(ctx context.Context, postID int64) error
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) CreatePost(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (user_id, title, content, image_path) VALUES ($1, $2, $3, $4) RETURNING post_id`
	return pr.db.QueryRowContext(ctx, query, post.UserID, post.Title, post.Content, post.ImagePath).Scan(&post.PostID)
}

func (pr *postRepository) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT post_id, user_id, title, content, created_at, image_path FROM posts WHERE post_id = $1`
	var post models.Post
	if err := pr.db.GetContext(ctx, &post, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (pr *postRepository) GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	query := `SELECT post_id, user_id, title, content, created_at, image_path FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var posts []models.Post
	if err := pr.db.SelectContext(ctx, &posts, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *postRepository) DeletePostByID(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE post_id = $1`
	result, err := pr.db.ExecContext(ctx, query, postID)
	return CheckErrResult(result, err, ErrPostNotFound)
}
