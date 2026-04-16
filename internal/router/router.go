package router

import (
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/handler"
	"socialnet/internal/middleware"
)

type Router struct {
	mux              *http.ServeMux
	userHandler      handler.UserHandler
	followingHandler handler.FollowingHandler
	messageHandler   handler.MessageHandler
	postHandler      handler.PostHandler
	uploadHandler    handler.UploadHandler
	commentHandler   handler.CommentHandler
	bookmarkHandler  handler.BookmarkHandler
	likeHandler      handler.LikeHandler
	authenticationMw *middleware.AuthenticationMiddleware
	rateLimit        *middleware.RateLimitMiddleware
}

func NewRouter(db *database.Database, cache *cache.Cache) *Router {
	return &Router{
		mux:              http.NewServeMux(),
		authenticationMw: middleware.NewAuthenticationMiddleware(db.SessionRepo),
		rateLimit:        middleware.NewRateLimitMiddleware(),
		commentHandler:   handler.NewCommentHandler(db.CommentRepo, cache),
		postHandler:      handler.NewPostHandler(db.PostRepo, cache),
		bookmarkHandler:  handler.NewBookmarkHandler(db.BookmarkRepo, cache),
		userHandler:      handler.NewUserHandler(db.UserRepo, db.SessionRepo, cache),
		likeHandler:      handler.NewLikeHandler(db.LikeRepo, cache),
		messageHandler:   handler.NewMessageHandler(db.MessageRepo, cache),
		uploadHandler:    handler.NewUploadHandler(),
		followingHandler: handler.NewFollowingHandler(db.FollowRepo, cache),
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	//session
	r.mux.HandleFunc("POST /api/v1/auth/register", r.rateLimit.RateLimit(r.userHandler.RegisterUser))
	r.mux.HandleFunc("POST /api/v1/auth/login", r.rateLimit.RateLimit(r.userHandler.LoginUser))
	r.mux.HandleFunc("DELETE /api/v1/auth/logout", r.authenticationMw.AuthenticationMiddleware(r.userHandler.LogoutUser))

	//user - me
	r.mux.HandleFunc("GET /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.GetCurrentUser))
	r.mux.HandleFunc("PUT /api/v1/users/me/info", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateUserInfo))
	r.mux.HandleFunc("PUT /api/v1/users/me/password", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateUserPassword))
	r.mux.HandleFunc("PUT /api/v1/users/me/profile_picture", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateProfilePicture))
	r.mux.HandleFunc("PUT /api/v1/users/me/banner", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateUserBanner))
	r.mux.HandleFunc("DELETE /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.DeleteCurrentUser))

	r.mux.HandleFunc("GET /api/v1/users/{user_id}", r.userHandler.GetUserByID)

	r.mux.HandleFunc("POST /api/v1/upload", r.rateLimit.RateLimit(r.authenticationMw.AuthenticationMiddleware(r.uploadHandler.Upload)))

	//following
	r.mux.HandleFunc("POST /api/v1/following/follow/{user_id}", r.authenticationMw.AuthenticationMiddleware(r.followingHandler.FollowUser))
	r.mux.HandleFunc("DELETE /api/v1/following/unfollow/{user_id}", r.authenticationMw.AuthenticationMiddleware(r.followingHandler.UnfollowUser))
	r.mux.HandleFunc("GET /api/v1/users/{user_id}/followers", r.followingHandler.GetFollowers)
	r.mux.HandleFunc("GET /api/v1/users/{user_id}/followers/count", r.followingHandler.GetCantFollowers)
	r.mux.HandleFunc("GET /api/v1/users/{user_id}/following", r.followingHandler.GetFollowing)
	r.mux.HandleFunc("GET /api/v1/users/{user_id}/following/count", r.followingHandler.GetCantFollowing)

	//message
	r.mux.HandleFunc("POST /api/v1/chats", r.authenticationMw.AuthenticationMiddleware(r.messageHandler.CreateChat))
	r.mux.HandleFunc("GET /api/v1/chats/me", r.authenticationMw.AuthenticationMiddleware(r.messageHandler.GetChatsFromUser))

	r.mux.HandleFunc("DELETE /api/v1/chats/messages/{message_id}", r.authenticationMw.AuthenticationMiddleware(r.messageHandler.DeleteMessage))

	r.mux.HandleFunc("GET /api/v1/chats/{chat_id}/messages", r.authenticationMw.AuthenticationMiddleware(r.messageHandler.GetMessagesFromChat))
	r.mux.HandleFunc("POST /api/v1/chats/{chat_id}/messages", r.authenticationMw.AuthenticationMiddleware(r.messageHandler.CreateMessage))

	//post
	r.mux.HandleFunc("POST /api/v1/posts", r.rateLimit.RateLimit(r.authenticationMw.AuthenticationMiddleware(r.postHandler.CreatePost)))
	r.mux.HandleFunc("GET /api/v1/posts/{post_id}", r.postHandler.GetPost)
	r.mux.HandleFunc("DELETE /api/v1/posts/{post_id}", r.authenticationMw.AuthenticationMiddleware(r.postHandler.DeletePost))

	r.mux.HandleFunc("GET /api/v1/users/{user_id}/posts", r.postHandler.GetPostsFromUser)
	r.mux.HandleFunc("GET /api/v1/feed", r.rateLimit.RateLimit(r.authenticationMw.AuthenticationMiddleware(r.postHandler.GetPostsFromFollowing)))

	//comments
	r.mux.HandleFunc("POST /api/v1/posts/{post_id}/comments", r.authenticationMw.AuthenticationMiddleware(r.commentHandler.CreateCommentForPost))
	r.mux.HandleFunc("POST /api/v1/comments/{comment_id}/comments", r.authenticationMw.AuthenticationMiddleware(r.commentHandler.CreateCommentForComment))
	r.mux.HandleFunc("GET /api/v1/comments/{comment_id}", r.commentHandler.GetComment)
	r.mux.HandleFunc("GET /api/v1/posts/{post_id}/comments", r.commentHandler.GetCommentsFromPost)
	r.mux.HandleFunc("GET /api/v1/comments/{comment_id}/comments", r.commentHandler.GetCommentsFromComment)
	r.mux.HandleFunc("DELETE /api/v1/comments/{comment_id}", r.authenticationMw.AuthenticationMiddleware(r.commentHandler.DeleteComment))
	r.mux.HandleFunc("GET /api/v1/users/me/comments", r.authenticationMw.AuthenticationMiddleware(r.commentHandler.GetCommentsFromUser))

	//bookmark
	r.mux.HandleFunc("POST /api/v1/posts/{post_id}/bookmark", r.authenticationMw.AuthenticationMiddleware(r.bookmarkHandler.CreateBookmark))
	r.mux.HandleFunc("DELETE /api/v1/posts/{post_id}/bookmark", r.authenticationMw.AuthenticationMiddleware(r.bookmarkHandler.DeleteBookmark))
	r.mux.HandleFunc("GET /api/v1/users/me/bookmarks", r.authenticationMw.AuthenticationMiddleware(r.bookmarkHandler.GetPostsFromUsersBookmark))

	//like
	r.mux.HandleFunc("POST /api/v1/posts/{post_id}/like", r.authenticationMw.AuthenticationMiddleware(r.likeHandler.CreateLikeForPost))
	r.mux.HandleFunc("POST /api/v1/comments/{comment_id}/like", r.authenticationMw.AuthenticationMiddleware(r.likeHandler.CreateLikeForComment))

	r.mux.HandleFunc("DELETE /api/v1/posts/{post_id}/like", r.authenticationMw.AuthenticationMiddleware(r.likeHandler.DeleteLikeForPost))
	r.mux.HandleFunc("DELETE /api/v1/comments/{comment_id}/like", r.authenticationMw.AuthenticationMiddleware(r.likeHandler.DeleteLikeForComment))

	r.mux.HandleFunc("GET /api/v1/users/me/likes", r.authenticationMw.AuthenticationMiddleware(r.likeHandler.GetPostsFromUsersLikes))
	r.mux.HandleFunc("GET /api/v1/posts/{post_id}/likes/count", r.likeHandler.GetLikesCountFromPost)
	r.mux.HandleFunc("GET /api/v1/comments/{comment_id}/likes/count", r.likeHandler.GetLikesCountFromComment)

	return r.mux
}
