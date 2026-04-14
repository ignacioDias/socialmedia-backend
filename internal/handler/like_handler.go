package handler

import (
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/models"
	"socialnet/internal/service"
	"strconv"
)

type LikeHandler interface {
	CreateLikeForPost(w http.ResponseWriter, r *http.Request)
	CreateLikeForComment(w http.ResponseWriter, r *http.Request)
	GetPostsFromUsersLikes(w http.ResponseWriter, r *http.Request)
	GetLikesCountFromPost(w http.ResponseWriter, r *http.Request)
	GetLikesCountFromComment(w http.ResponseWriter, r *http.Request)
	DeleteLikeForPost(w http.ResponseWriter, r *http.Request)
	DeleteLikeForComment(w http.ResponseWriter, r *http.Request)
}

type likeHandler struct {
	likeService service.LikeService
}

func NewLikeHandler(likeRepo database.LikeRepository, cache *cache.Cache) LikeHandler {
	return &likeHandler{
		likeService: service.NewLikeService(likeRepo, cache),
	}
}

func (h *likeHandler) CreateLikeForPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idFromPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	like := models.NewLike(postID, models.PostTarget, userID)
	if err := h.likeService.CreateLike(r.Context(), like); err != nil {
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, like, http.StatusCreated)
}

func (h *likeHandler) CreateLikeForComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idFromPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	like := models.NewLike(commentID, models.CommentTarget, userID)
	if err := h.likeService.CreateLike(r.Context(), like); err != nil {
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, like, http.StatusCreated)
}

func (h *likeHandler) GetPostsFromUsersLikes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	limit, err := ParseQueryInt(r, "limit", DEFAULT_LIMIT)
	if err != nil {
		http.Error(w, "invalid format for limit ", http.StatusBadRequest)
		return
	}
	offset, err := ParseQueryInt(r, "offset", DEFAULT_OFFSET)
	if err != nil {
		http.Error(w, "invalid format for offset ", http.StatusBadRequest)
		return
	}
	posts, err := h.likeService.GetPostsFromUsersLikes(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get posts from users likes", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, posts, http.StatusOK)
}

func (h *likeHandler) GetLikesCountFromPost(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	count, err := h.likeService.GetLikesCountFromTarget(r.Context(), postID, models.PostTarget)
	if err != nil {
		http.Error(w, "Failed to get likes count from post", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, map[string]int64{"likes_count": count}, http.StatusOK)
}

func (h *likeHandler) GetLikesCountFromComment(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	count, err := h.likeService.GetLikesCountFromTarget(r.Context(), commentID, models.CommentTarget)
	if err != nil {
		http.Error(w, "Failed to get likes count from comment", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, map[string]int64{"likes_count": count}, http.StatusOK)
}

func (h *likeHandler) DeleteLikeForPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idFromPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	like := models.NewLike(postID, models.PostTarget, userID)
	if err := h.likeService.DeleteLike(r.Context(), like); err != nil {
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *likeHandler) DeleteLikeForComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idFromPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	like := models.NewLike(commentID, models.CommentTarget, userID)
	if err := h.likeService.DeleteLike(r.Context(), like); err != nil {
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
