package handler

import (
	"encoding/json"
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/models"
	"socialnet/internal/service"
	"strconv"
)

type PostHandler interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	GetPostsFromUser(w http.ResponseWriter, r *http.Request)
	GetPostsFromFollowing(w http.ResponseWriter, r *http.Request)
}

type implPostHandler struct {
	postService service.PostService
}

func NewPostHandler(postRepo database.PostRepository, cache *cache.Cache) PostHandler {
	return &implPostHandler{
		postService: service.NewPostService(postRepo, cache),
	}
}

func (ph *implPostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var postReq service.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		http.Error(w, "wrong format for body", http.StatusBadRequest)
		return
	}
	post := models.NewPost(userID, postReq.Title, postReq.ImagePath, postReq.Content)
	if err := ph.postService.CreatePost(r.Context(), post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, post, http.StatusCreated)
}

func (ph *implPostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for post_id", http.StatusBadRequest)
		return
	}
	post, err := ph.postService.GetPostFromID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, post, http.StatusOK)
}

func (ph *implPostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idFromPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for post_id", http.StatusBadRequest)
		return
	}
	if err := ph.postService.DeletePost(r.Context(), postID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ph *implPostHandler) GetPostsFromUser(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("user_id")
	userID, err := strconv.ParseInt(idFromPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for post_id", http.StatusBadRequest)
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
	posts, err := ph.postService.GetPostsFromUser(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, posts, http.StatusOK)
}

func (ph *implPostHandler) GetPostsFromFollowing(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
	posts, err := ph.postService.GetPostsFromFollowing(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, posts, http.StatusOK)
}
