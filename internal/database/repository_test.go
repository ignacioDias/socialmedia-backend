package database

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"socialnet/internal/models"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	rawDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db := sqlx.NewDb(rawDB, "sqlmock")
	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet SQL expectations: %v", err)
		}
		_ = rawDB.Close()
	}

	return db, mock, cleanup
}

func TestCheckErrResult(t *testing.T) {
	errBoom := errors.New("boom")
	errNotFound := errors.New("not found")

	if err := CheckErrResult(sqlmock.NewResult(0, 1), nil, errNotFound); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if err := CheckErrResult(sqlmock.NewResult(0, 0), nil, errNotFound); !errors.Is(err, errNotFound) {
		t.Fatalf("expected not-found error, got: %v", err)
	}

	if err := CheckErrResult(sqlmock.NewResult(0, 1), errBoom, errNotFound); !errors.Is(err, errBoom) {
		t.Fatalf("expected original error, got: %v", err)
	}
}

func TestPostRepository_CreatePost(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewPostRepository(db)
	post := &models.Post{UserID: 1, Title: "hello", Content: "world", ImagePath: "img.png"}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO posts (user_id, title, content, image_path) VALUES ($1, $2, $3, $4) RETURNING post_id")).
		WithArgs(post.UserID, post.Title, post.Content, post.ImagePath).
		WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(42)))

	if err := repo.CreatePost(context.Background(), post); err != nil {
		t.Fatalf("CreatePost returned error: %v", err)
	}
	if post.PostID != 42 {
		t.Fatalf("expected post id 42, got: %d", post.PostID)
	}
}

func TestPostRepository_GetPostByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewPostRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT post_id, user_id, title, content, created_at, image_path FROM posts WHERE post_id = $1")).
		WithArgs(int64(777)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetPostByID(context.Background(), 777)
	if !errors.Is(err, ErrPostNotFound) {
		t.Fatalf("expected ErrPostNotFound, got: %v", err)
	}
}

func TestCommentRepository_CreateComment(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewCommentRepository(db)
	comment := &models.Comment{TargetID: 9, TargetType: models.PostTarget, UserID: 5, Content: "nice", ImagePath: ""}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO comments (target_id, target_type, user_id, content, image_path) VALUES ($1, $2, $3, $4, $5) RETURNING comment_id")).
		WithArgs(comment.TargetID, comment.TargetType, comment.UserID, comment.Content, comment.ImagePath).
		WillReturnRows(sqlmock.NewRows([]string{"comment_id"}).AddRow(int64(12)))

	if err := repo.CreateComment(context.Background(), comment); err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}
	if comment.CommentID != 12 {
		t.Fatalf("expected comment id 12, got: %d", comment.CommentID)
	}
}

func TestCommentRepository_GetCommentByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewCommentRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT comment_id, target_id, target_type, user_id, content, created_at FROM comments WHERE comment_id = $1")).
		WithArgs(int64(444)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetCommentByID(context.Background(), 444)
	if !errors.Is(err, ErrCommentNotFound) {
		t.Fatalf("expected ErrCommentNotFound, got: %v", err)
	}
}

func TestFollowRepository_DeleteFollow_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewFollowRepository(db)
	follow := &models.Follow{FollowerID: 1, FollowingID: 2}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM follows WHERE follower_id = $1 AND following_id = $2")).
		WithArgs(follow.FollowerID, follow.FollowingID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteFollow(context.Background(), follow)
	if !errors.Is(err, ErrFollowNotFound) {
		t.Fatalf("expected ErrFollowNotFound, got: %v", err)
	}
}

func TestLikeRepository_GetLikesCountFromTarget(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewLikeRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM likes WHERE target_id = $1 AND target_type = $2")).
		WithArgs(int64(55), models.PostTarget).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(8)))

	count, err := repo.GetLikesCountFromTarget(context.Background(), 55, models.PostTarget)
	if err != nil {
		t.Fatalf("GetLikesCountFromTarget returned error: %v", err)
	}
	if count != 8 {
		t.Fatalf("expected likes count 8, got: %d", count)
	}
}

func TestLikeRepository_DeleteLike_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewLikeRepository(db)
	like := &models.Like{TargetType: models.CommentTarget, TargetID: 3, UserID: 9}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM likes WHERE target_type = $1 AND target_id = $2 AND user_id = $3")).
		WithArgs(like.TargetType, like.TargetID, like.UserID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteLike(context.Background(), like)
	if !errors.Is(err, ErrLikeNotFound) {
		t.Fatalf("expected ErrLikeNotFound, got: %v", err)
	}
}

func TestBookmarkRepository_CreateBookmark(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewBookmarkRepository(db)
	bookmark := &models.Bookmark{PostID: 7, UserID: 11}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO bookmarks (post_id, user_id) VALUES ($1, $2)")).
		WithArgs(bookmark.PostID, bookmark.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.CreateBookmark(context.Background(), bookmark); err != nil {
		t.Fatalf("CreateBookmark returned error: %v", err)
	}
}

func TestRepostRepository_DeleteRepost_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewRepostRepository(db)
	repost := &models.Repost{UserID: 3, PostID: 10}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM reposts WHERE user_id = $1 AND post_id = $2")).
		WithArgs(repost.UserID, repost.PostID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteRepost(context.Background(), repost)
	if !errors.Is(err, ErrRepostNotFound) {
		t.Fatalf("expected ErrRepostNotFound, got: %v", err)
	}
}

func TestRepostRepository_GetRepostsCountFromPost(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewRepostRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM reposts WHERE post_id = $1")).
		WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(4)))

	count, err := repo.GetRepostsCountFromPost(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetRepostsCountFromPost returned error: %v", err)
	}
	if count != 4 {
		t.Fatalf("expected count 4, got: %d", count)
	}
}

func TestMessageRepository_CreateChat(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	chat := &models.Chat{User1ID: 1, User2ID: 2}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO chats (user1_id, user2_id) VALUES ($1, $2) RETURNING chat_id")).
		WithArgs(chat.User1ID, chat.User2ID).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id"}).AddRow(int64(99)))

	if err := repo.CreateChat(context.Background(), chat); err != nil {
		t.Fatalf("CreateChat returned error: %v", err)
	}
	if chat.ChatID != 99 {
		t.Fatalf("expected chat id 99, got: %d", chat.ChatID)
	}
}

func TestMessageRepository_DeleteMessageByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM messages WHERE message_id = $1 AND sender_id = $2")).
		WithArgs(int64(123), int64(45)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteMessageByID(context.Background(), 123, 45)
	if !errors.Is(err, ErrMessageNotFound) {
		t.Fatalf("expected ErrMessageNotFound, got: %v", err)
	}
}

func TestSessionRepository_FindSessionByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewSessionRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1 AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')")).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.FindSessionByID(context.Background(), "missing")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
}

func TestSessionRepository_DeleteSessionByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewSessionRepository(db)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sessions WHERE id = $1")).
		WithArgs("abc").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteSessionByID(context.Background(), "abc")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got: %v", err)
	}
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).
		WithArgs("none@example.com").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetUserByEmail(context.Background(), "none@example.com")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_UpdateUserInfo_Success(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	user := &models.User{UserID: 3, Username: "new_name", Email: "new@example.com"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET username = $1, email = $2 WHERE user_id = $3")).
		WithArgs(user.Username, user.Email, user.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updatedUser, err := repo.UpdateUserInfo(context.Background(), user)
	if err != nil {
		t.Fatalf("UpdateUserInfo returned error: %v", err)
	}
	if updatedUser.UserID != user.UserID {
		t.Fatalf("unexpected returned user id: %d", updatedUser.UserID)
	}
}

func TestMessageRepository_GetChatsFromUserID(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewMessageRepository(db)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM chats WHERE (user1_id = $1 OR user2_id = $1) ORDER BY last_message_at DESC")).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user1_id", "user2_id", "last_message_at"}).
			AddRow(int64(1), int64(5), int64(9), now).
			AddRow(int64(2), int64(4), int64(5), now))

	chats, err := repo.GetChatsFromUserID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetChatsFromUserID returned error: %v", err)
	}
	if len(chats) != 2 {
		t.Fatalf("expected 2 chats, got: %d", len(chats))
	}
}
