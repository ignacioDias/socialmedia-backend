package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db          *sqlx.DB
	PostRepo    PostRepository
	RepostRepo  RepostRepository
	SessionRepo SessionRepository
	CommentRepo CommentRepository
	LikeRepo    LikeRepository
}

var createPostsTable = `
CREATE TABLE IF NOT EXISTS posts (
	post_id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	image_path TEXT,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);`

var createCommentsTable = `
CREATE TABLE IF NOT EXISTS comments(
	comment_id BIGSERIAL PRIMARY KEY,
	target_id BIGINT NOT NULL,
	image_path TEXT,
	target_type TEXT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);`

var createLikesTable = `
CREATE TABLE IF NOT EXISTS likes(
	target_id BIGINT NOT NULL,
	target_type TEXT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW()
    PRIMARY KEY (target_id, target_type, user_id)
)`

var createRepostsTable = `
CREATE TABLE IF NOT EXISTS reposts(
	post_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	content TEXT,
	image_path TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW()
    PRIMARY KEY (post_id, user_id)
)`

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db:          db,
		PostRepo:    NewPostRepository(db),
		LikeRepo:    NewLikeRepository(db),
		SessionRepo: NewSessionRepository(db),
		RepostRepo:  NewRepostRepository(db),
		CommentRepo: NewCommentRepository(db),
	}
}

func (db *Database) Init() error {
	tx, err := db.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin init db transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	tableStatements := []struct {
		name string
		ddl  string
	}{
		{name: "posts", ddl: createPostsTable},
		{name: "likes", ddl: createLikesTable},
		{name: "reposts", ddl: createRepostsTable},
		{name: "comments", ddl: createCommentsTable},
	}

	for _, tableStmt := range tableStatements {
		if _, err := tx.Exec(tableStmt.ddl); err != nil {
			return fmt.Errorf("create table %s: %w", tableStmt.name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit init db transaction: %w", err)
	}

	return nil
}
