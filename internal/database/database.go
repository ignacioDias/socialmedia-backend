package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db           *sqlx.DB
	PostRepo     PostRepository
	FollowRepo   FollowRepository
	RepostRepo   RepostRepository
	SessionRepo  SessionRepository
	CommentRepo  CommentRepository
	LikeRepo     LikeRepository
	BookmarkRepo BookmarkRepository
	MessageRepo  MessageRepository
	UserRepo     UserRepository
}

var createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);`

var createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	user_id BIGSERIAL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	hashed_password TEXT NOT NULL,
	profile_picture_path TEXT NOT NULL, 
	banner_path TEXT NOT NULL
);`
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
	created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (target_id, target_type, user_id)
);`

var createFollowsTable = `
CREATE TABLE IF NOT EXISTS follows(
	following_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	follower_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY (follower_id, following_id)
);`
var createRepostsTable = `
CREATE TABLE IF NOT EXISTS reposts(
	post_id BIGINT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	content TEXT,
	image_path TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);`
var createBookmarksTable = `
CREATE TABLE IF NOT EXISTS bookmarks(
	post_id BIGINT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
	user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY (user_id, post_id)
);`

var createChatsTable = `
CREATE TABLE IF NOT EXISTS chats(
    chat_id BIGSERIAL PRIMARY KEY,
    user1_id BIGINT NOT NULL REFERENCES users(user_id),
    user2_id BIGINT NOT NULL REFERENCES users(user_id),
    last_message_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT ordered_users CHECK (user1_id < user2_id),
    CONSTRAINT unique_chat UNIQUE (user1_id, user2_id)
);`
var createMessagesTable = `
CREATE TABLE IF NOT EXISTS messages(
	message_id BIGSERIAL PRIMARY KEY,
	chat_id BIGINT NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
	sender_id BIGINT NOT NULL REFERENCES users(user_id),
	image_path TEXT,
	content TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW()
);`

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db:           db,
		PostRepo:     NewPostRepository(db),
		LikeRepo:     NewLikeRepository(db),
		SessionRepo:  NewSessionRepository(db),
		BookmarkRepo: NewBookmarkRepository(db),
		RepostRepo:   NewRepostRepository(db),
		CommentRepo:  NewCommentRepository(db),
		FollowRepo:   NewFollowRepository(db),
		MessageRepo:  NewMessageRepository(db),
		UserRepo:     NewUserRepository(db),
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
		{name: "users", ddl: createUsersTable},
		{name: "sessions", ddl: createSessionsTable},
		{name: "posts", ddl: createPostsTable},
		{name: "likes", ddl: createLikesTable},
		{name: "bookmarks", ddl: createBookmarksTable},
		{name: "reposts", ddl: createRepostsTable},
		{name: "comments", ddl: createCommentsTable},
		{name: "chats", ddl: createChatsTable},
		{name: "follows", ddl: createFollowsTable},
		{name: "messages", ddl: createMessagesTable},
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
