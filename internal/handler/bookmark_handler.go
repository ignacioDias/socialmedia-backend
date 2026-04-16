package handler

import (
	"errors"
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/models"
	"socialnet/internal/service"
	"strconv"
)

type BookmarkHandler interface {
	CreateBookmark(w http.ResponseWriter, r *http.Request)
	DeleteBookmark(w http.ResponseWriter, r *http.Request)
	GetPostsFromUsersBookmark(w http.ResponseWriter, r *http.Request)
}

type implBookmarkHandler struct {
	bookmarkService service.BookmarkService
}

func NewBookmarkHandler(bookmarkRepo database.BookmarkRepository, cache *cache.Cache) BookmarkHandler {
	return &implBookmarkHandler{
		bookmarkService: service.NewBookmarkService(bookmarkRepo, cache),
	}
}

func (bh *implBookmarkHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	idPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
		return
	}
	bookmark := models.NewBookmark(postID, userID)
	if err := bh.bookmarkService.CreateBookmark(r.Context(), bookmark); err != nil {
		http.Error(w, "error creating bookmark", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, bookmark, http.StatusCreated)
}

func (bh *implBookmarkHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	idPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
		return
	}
	bookmark := models.NewBookmark(postID, userID)
	if err := bh.bookmarkService.DeleteBookmark(r.Context(), bookmark); err != nil {
		if errors.Is(err, database.ErrBookmarkNotFound) {
			http.Error(w, "bookmark not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting bookmark", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (bh *implBookmarkHandler) GetPostsFromUsersBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
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
	posts, err := bh.bookmarkService.GetPostsFromUsersBookmarks(r.Context(), userID, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		http.Error(w, "error fetching bookmarked posts", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, posts, http.StatusOK)
}
