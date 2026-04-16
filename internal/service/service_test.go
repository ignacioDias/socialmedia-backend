package service

import (
	"context"
	"errors"
	"testing"

	"socialnet/internal/cache"
	"socialnet/internal/models"
)

func newTestCache() *cache.Cache {
	return cache.NewCache("127.0.0.1:0")
}

type postRepoStub struct {
	createPostFn          func(ctx context.Context, post *models.Post) error
	getPostByIDFn         func(ctx context.Context, id int64) (*models.Post, error)
	getPostsByUserIDFn    func(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	getPostsFromFollowsFn func(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	deletePostByIDFn      func(ctx context.Context, postID, userID int64) error
}

func (s *postRepoStub) CreatePost(ctx context.Context, post *models.Post) error {
	if s.createPostFn != nil {
		return s.createPostFn(ctx, post)
	}
	return nil
}
func (s *postRepoStub) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	if s.getPostByIDFn != nil {
		return s.getPostByIDFn(ctx, id)
	}
	return nil, nil
}
func (s *postRepoStub) GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if s.getPostsByUserIDFn != nil {
		return s.getPostsByUserIDFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *postRepoStub) GetPostsFromFollows(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if s.getPostsFromFollowsFn != nil {
		return s.getPostsFromFollowsFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *postRepoStub) DeletePostByID(ctx context.Context, postID, userID int64) error {
	if s.deletePostByIDFn != nil {
		return s.deletePostByIDFn(ctx, postID, userID)
	}
	return nil
}

type commentRepoStub struct {
	createCommentFn         func(ctx context.Context, comment *models.Comment) error
	getCommentByIDFn        func(ctx context.Context, commentID int64) (*models.Comment, error)
	getCommentsFromTargetFn func(ctx context.Context, targetType models.TargetType, targetID int64, limit, offset int) ([]models.Comment, error)
	getCommentsFromUserFn   func(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error)
	deleteCommentByIDFn     func(ctx context.Context, commentID, userID int64) error
}

func (s *commentRepoStub) CreateComment(ctx context.Context, comment *models.Comment) error {
	if s.createCommentFn != nil {
		return s.createCommentFn(ctx, comment)
	}
	return nil
}
func (s *commentRepoStub) GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error) {
	if s.getCommentByIDFn != nil {
		return s.getCommentByIDFn(ctx, commentID)
	}
	return nil, nil
}
func (s *commentRepoStub) GetCommentsFromTarget(ctx context.Context, targetType models.TargetType, targetID int64, limit, offset int) ([]models.Comment, error) {
	if s.getCommentsFromTargetFn != nil {
		return s.getCommentsFromTargetFn(ctx, targetType, targetID, limit, offset)
	}
	return nil, nil
}
func (s *commentRepoStub) GetCommentsFromUser(ctx context.Context, userID int64, limit, offset int) ([]models.Comment, error) {
	if s.getCommentsFromUserFn != nil {
		return s.getCommentsFromUserFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *commentRepoStub) DeleteCommentByID(ctx context.Context, commentID, userID int64) error {
	if s.deleteCommentByIDFn != nil {
		return s.deleteCommentByIDFn(ctx, commentID, userID)
	}
	return nil
}

type followRepoStub struct {
	createFollowFn              func(ctx context.Context, follow *models.Follow) error
	getFollowersFromUserIDFn    func(ctx context.Context, userID int64, limit, offset int) ([]models.User, error)
	getFollowingFromUserIDFn    func(ctx context.Context, userID int64, limit, offset int) ([]models.User, error)
	getCantFollowersFromUserID  func(ctx context.Context, userID int64) (int, error)
	getCantFollowingsFromUserID func(ctx context.Context, userID int64) (int, error)
	deleteFollowFn              func(ctx context.Context, follow *models.Follow) error
}

func (s *followRepoStub) CreateFollow(ctx context.Context, follow *models.Follow) error {
	if s.createFollowFn != nil {
		return s.createFollowFn(ctx, follow)
	}
	return nil
}
func (s *followRepoStub) GetFollowersFromUserID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	if s.getFollowersFromUserIDFn != nil {
		return s.getFollowersFromUserIDFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *followRepoStub) GetFollowingFromUserID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	if s.getFollowingFromUserIDFn != nil {
		return s.getFollowingFromUserIDFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *followRepoStub) GetCantFollowersFromUserID(ctx context.Context, userID int64) (int, error) {
	if s.getCantFollowersFromUserID != nil {
		return s.getCantFollowersFromUserID(ctx, userID)
	}
	return 0, nil
}
func (s *followRepoStub) GetCantFollowingsFromUserID(ctx context.Context, userID int64) (int, error) {
	if s.getCantFollowingsFromUserID != nil {
		return s.getCantFollowingsFromUserID(ctx, userID)
	}
	return 0, nil
}
func (s *followRepoStub) DeleteFollow(ctx context.Context, follow *models.Follow) error {
	if s.deleteFollowFn != nil {
		return s.deleteFollowFn(ctx, follow)
	}
	return nil
}

type likeRepoStub struct {
	createLikeFn             func(ctx context.Context, like *models.Like) error
	getPostsFromUsersLikesFn func(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
	getLikesCountFromTarget  func(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error)
	deleteLikeFn             func(ctx context.Context, like *models.Like) error
}

func (s *likeRepoStub) CreateLike(ctx context.Context, like *models.Like) error {
	if s.createLikeFn != nil {
		return s.createLikeFn(ctx, like)
	}
	return nil
}
func (s *likeRepoStub) GetPostsFromUsersLikes(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if s.getPostsFromUsersLikesFn != nil {
		return s.getPostsFromUsersLikesFn(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (s *likeRepoStub) GetLikesCountFromTarget(ctx context.Context, targetID int64, targetType models.TargetType) (int64, error) {
	if s.getLikesCountFromTarget != nil {
		return s.getLikesCountFromTarget(ctx, targetID, targetType)
	}
	return 0, nil
}
func (s *likeRepoStub) DeleteLike(ctx context.Context, like *models.Like) error {
	if s.deleteLikeFn != nil {
		return s.deleteLikeFn(ctx, like)
	}
	return nil
}

type messageRepoStub struct {
	createMessageFn      func(ctx context.Context, message *models.Message) error
	deleteMessageByIDFn  func(ctx context.Context, messageID, userID int64) error
	createChatFn         func(ctx context.Context, chat *models.Chat) error
	getChatsFromUserIDFn func(ctx context.Context, userID int64) ([]models.Chat, error)
	getMessagesFromChat  func(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error)
}

func (s *messageRepoStub) CreateMessage(ctx context.Context, message *models.Message) error {
	if s.createMessageFn != nil {
		return s.createMessageFn(ctx, message)
	}
	return nil
}
func (s *messageRepoStub) DeleteMessageByID(ctx context.Context, messageID, userID int64) error {
	if s.deleteMessageByIDFn != nil {
		return s.deleteMessageByIDFn(ctx, messageID, userID)
	}
	return nil
}
func (s *messageRepoStub) CreateChat(ctx context.Context, chat *models.Chat) error {
	if s.createChatFn != nil {
		return s.createChatFn(ctx, chat)
	}
	return nil
}
func (s *messageRepoStub) GetChatsFromUserID(ctx context.Context, userID int64) ([]models.Chat, error) {
	if s.getChatsFromUserIDFn != nil {
		return s.getChatsFromUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (s *messageRepoStub) GetMessagesFromChat(ctx context.Context, chatID, userID int64, limit, offset int) ([]models.Message, error) {
	if s.getMessagesFromChat != nil {
		return s.getMessagesFromChat(ctx, chatID, userID, limit, offset)
	}
	return nil, nil
}

type bookmarkRepoStub struct {
	createBookmarkFn          func(ctx context.Context, bookmark *models.Bookmark) error
	deleteBookmarkFn          func(ctx context.Context, bookmark *models.Bookmark) error
	getPostsFromUsersBookmark func(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error)
}

func (s *bookmarkRepoStub) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	if s.createBookmarkFn != nil {
		return s.createBookmarkFn(ctx, bookmark)
	}
	return nil
}
func (s *bookmarkRepoStub) DeleteBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	if s.deleteBookmarkFn != nil {
		return s.deleteBookmarkFn(ctx, bookmark)
	}
	return nil
}
func (s *bookmarkRepoStub) GetPostsFromUsersBookmark(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	if s.getPostsFromUsersBookmark != nil {
		return s.getPostsFromUsersBookmark(ctx, userID, limit, offset)
	}
	return nil, nil
}

type userRepoStub struct {
	createUserFn               func(ctx context.Context, user *models.User) error
	getUserByIDFn              func(ctx context.Context, userID int64) (*models.User, error)
	getUserByEmailFn           func(ctx context.Context, email string) (*models.User, error)
	getUserByUsernameFn        func(ctx context.Context, username string) (*models.User, error)
	updateUserInfoFn           func(ctx context.Context, user *models.User) (*models.User, error)
	updateUserProfilePictureFn func(ctx context.Context, profilePicturePath string, userID int64) error
	updateUserBannerFn         func(ctx context.Context, bannerPath string, userID int64) error
	updateUserHashedPasswordFn func(ctx context.Context, hashedPassword string, userID int64) error
	deleteUserByIDFn           func(ctx context.Context, userID int64) error
}

func (s *userRepoStub) CreateUser(ctx context.Context, user *models.User) error {
	if s.createUserFn != nil {
		return s.createUserFn(ctx, user)
	}
	return nil
}
func (s *userRepoStub) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	if s.getUserByIDFn != nil {
		return s.getUserByIDFn(ctx, userID)
	}
	return nil, nil
}
func (s *userRepoStub) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}
func (s *userRepoStub) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if s.getUserByUsernameFn != nil {
		return s.getUserByUsernameFn(ctx, username)
	}
	return nil, nil
}
func (s *userRepoStub) UpdateUserInfo(ctx context.Context, user *models.User) (*models.User, error) {
	if s.updateUserInfoFn != nil {
		return s.updateUserInfoFn(ctx, user)
	}
	return user, nil
}
func (s *userRepoStub) UpdateUserProfilePicture(ctx context.Context, profilePicturePath string, userID int64) error {
	if s.updateUserProfilePictureFn != nil {
		return s.updateUserProfilePictureFn(ctx, profilePicturePath, userID)
	}
	return nil
}
func (s *userRepoStub) UpdateUserBanner(ctx context.Context, bannerPath string, userID int64) error {
	if s.updateUserBannerFn != nil {
		return s.updateUserBannerFn(ctx, bannerPath, userID)
	}
	return nil
}
func (s *userRepoStub) UpdateUserHashedPassword(ctx context.Context, hashedPassword string, userID int64) error {
	if s.updateUserHashedPasswordFn != nil {
		return s.updateUserHashedPasswordFn(ctx, hashedPassword, userID)
	}
	return nil
}
func (s *userRepoStub) DeleteUserByID(ctx context.Context, userID int64) error {
	if s.deleteUserByIDFn != nil {
		return s.deleteUserByIDFn(ctx, userID)
	}
	return nil
}

type sessionRepoStub struct {
	createSessionFn          func(ctx context.Context, session *models.Session) error
	findSessionByIDFn        func(ctx context.Context, id string) (*models.Session, error)
	deleteSessionByIDFn      func(ctx context.Context, id string) error
	deleteSessionsByUserIDFn func(ctx context.Context, userID int64) error
}

func (s *sessionRepoStub) CreateSession(ctx context.Context, session *models.Session) error {
	if s.createSessionFn != nil {
		return s.createSessionFn(ctx, session)
	}
	return nil
}
func (s *sessionRepoStub) FindSessionByID(ctx context.Context, id string) (*models.Session, error) {
	if s.findSessionByIDFn != nil {
		return s.findSessionByIDFn(ctx, id)
	}
	return nil, nil
}
func (s *sessionRepoStub) DeleteSessionByID(ctx context.Context, id string) error {
	if s.deleteSessionByIDFn != nil {
		return s.deleteSessionByIDFn(ctx, id)
	}
	return nil
}
func (s *sessionRepoStub) DeleteSessionsByUserID(ctx context.Context, userID int64) error {
	if s.deleteSessionsByUserIDFn != nil {
		return s.deleteSessionsByUserIDFn(ctx, userID)
	}
	return nil
}

func TestPostService_GetPostFromID_InvalidID(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewPostService(&postRepoStub{}, c)
	_, err := svc.GetPostFromID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error for invalid id")
	}
}

func TestCommentService_DeleteCommentByID_InvalidID(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewCommentService(&commentRepoStub{}, c)
	err := svc.DeleteCommentByID(context.Background(), 0, 1)
	if err == nil {
		t.Fatalf("expected error for invalid comment id")
	}
}

func TestFollowingService_CreateFollow_CannotFollowSelf(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewFollowingService(&followRepoStub{}, c)
	err := svc.CreateFollow(context.Background(), &models.Follow{FollowerID: 10, FollowingID: 10})
	if err == nil {
		t.Fatalf("expected self-follow validation error")
	}
}

func TestFollowingService_GetCountFollowers_InvalidUserID(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewFollowingService(&followRepoStub{}, c)
	_, err := svc.GetCountFollowers(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected invalid user id error")
	}
}

func TestLikeService_CreateLike_WrapsRepoError(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	repoErr := errors.New("repo fail")
	svc := NewLikeService(&likeRepoStub{createLikeFn: func(ctx context.Context, like *models.Like) error {
		return repoErr
	}}, c)

	err := svc.CreateLike(context.Background(), &models.Like{TargetID: 1, TargetType: models.PostTarget, UserID: 2})
	if err == nil || !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got: %v", err)
	}
}

func TestLikeService_DeleteLike_InvalidLike(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewLikeService(&likeRepoStub{}, c)
	err := svc.DeleteLike(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected invalid like error")
	}
}

func TestMessageService_CreateChat_CannotChatWithSelf(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewMessageService(&messageRepoStub{}, c)
	err := svc.CreateChat(context.Background(), &models.Chat{User1ID: 3, User2ID: 3})
	if err == nil {
		t.Fatalf("expected self-chat validation error")
	}
}

func TestMessageService_CreateMessage_EmptyMessageRejected(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewMessageService(&messageRepoStub{}, c)
	err := svc.CreateMessage(context.Background(), &models.Message{Content: "", ImagePath: ""})
	if err == nil {
		t.Fatalf("expected empty message validation error")
	}
}

func TestBookmarkService_CreateBookmark_CallsRepo(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	called := false
	svc := NewBookmarkService(&bookmarkRepoStub{createBookmarkFn: func(ctx context.Context, bookmark *models.Bookmark) error {
		called = true
		return nil
	}}, c)

	err := svc.CreateBookmark(context.Background(), &models.Bookmark{PostID: 1, UserID: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected repository call")
	}
}

func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewUserService(&userRepoStub{}, &sessionRepoStub{}, c)
	_, err := svc.CreateUser(context.Background(), &UserRegisterRequest{
		Email:    "not-an-email",
		Username: "john",
		Password: "Aa1!aaaa",
	})
	if err == nil {
		t.Fatalf("expected invalid email error")
	}
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	h, err := hashPassword("Aa1!aaaa")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	svc := NewUserService(&userRepoStub{getUserByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
		return &models.User{UserID: 7, Email: email, HashedPassword: string(h)}, nil
	}}, &sessionRepoStub{}, c)

	_, err = svc.Login(context.Background(), &UserLoginRequest{Email: "john@example.com", Password: "wrong-pass"})
	if err == nil {
		t.Fatalf("expected invalid credentials")
	}
}

func TestUserService_Login_SuccessCreatesSession(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	h, err := hashPassword("Aa1!aaaa")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	sessionCreated := false
	svc := NewUserService(&userRepoStub{getUserByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
		return &models.User{UserID: 7, Email: email, HashedPassword: string(h)}, nil
	}}, &sessionRepoStub{createSessionFn: func(ctx context.Context, session *models.Session) error {
		sessionCreated = true
		return nil
	}}, c)

	sess, err := svc.Login(context.Background(), &UserLoginRequest{Email: "john@example.com", Password: "Aa1!aaaa"})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if sess == nil || sess.UserID != 7 {
		t.Fatalf("unexpected session: %+v", sess)
	}
	if !sessionCreated {
		t.Fatalf("expected CreateSession to be called")
	}
}

func TestUserService_UpdateInfo_InvalidEmail(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	svc := NewUserService(&userRepoStub{}, &sessionRepoStub{}, c)
	invalidEmail := "invalid"
	_, err := svc.UpdateInfo(context.Background(), 1, &UpdateInfoReq{Email: &invalidEmail})
	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got: %v", err)
	}
}

func TestUserService_Logout_Delegates(t *testing.T) {
	c := newTestCache()
	defer c.Close()

	called := false
	svc := NewUserService(&userRepoStub{}, &sessionRepoStub{deleteSessionsByUserIDFn: func(ctx context.Context, userID int64) error {
		called = true
		if userID != 9 {
			t.Fatalf("unexpected userID %d", userID)
		}
		return nil
	}}, c)

	if err := svc.Logout(context.Background(), 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected DeleteSessionsByUserID call")
	}
}
