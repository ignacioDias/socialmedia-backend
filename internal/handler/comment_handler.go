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

type CommentHandler interface {
	CreateCommentForPost(w http.ResponseWriter, r *http.Request)
	CreateCommentForComment(w http.ResponseWriter, r *http.Request)
	GetComment(w http.ResponseWriter, r *http.Request)
	GetCommentsFromPost(w http.ResponseWriter, r *http.Request)
	GetCommentsFromComment(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
	GetCommentsFromUser(w http.ResponseWriter, r *http.Request)
}

type implCommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentRepo database.CommentRepository, cache *cache.Cache) CommentHandler {
	return &implCommentHandler{
		commentService: service.NewCommentService(commentRepo, cache),
	}
}

func (ch *implCommentHandler) CreateCommentForPost(w http.ResponseWriter, r *http.Request) {
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
	var commentReq service.CommentReq
	if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
		http.Error(w, "wrong format for comment creation request", http.StatusBadRequest)
		return
	}
	comment := models.NewComment(postID, models.PostTarget, userID, commentReq.Content, commentReq.ImagePath)
	if err := ch.commentService.CreateComment(r.Context(), comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comment, http.StatusCreated)
}

func (ch *implCommentHandler) CreateCommentForComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	idPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
		return
	}
	var commentReq service.CommentReq
	if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
		http.Error(w, "wrong format for comment creation request", http.StatusBadRequest)
		return
	}
	comment := models.NewComment(commentID, models.CommentTarget, userID, commentReq.Content, commentReq.ImagePath)
	if err := ch.commentService.CreateComment(r.Context(), comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comment, http.StatusCreated)
}

func (ch *implCommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
		return
	}
	comment, err := ch.commentService.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comment, http.StatusOK)
}
func (ch *implCommentHandler) GetCommentsFromPost(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("post_id")
	postID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
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
	comments, err := ch.commentService.GetCommentsFromTarget(r.Context(), postID, models.PostTarget, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comments, http.StatusOK)
}
func (ch *implCommentHandler) GetCommentsFromComment(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
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
	comments, err := ch.commentService.GetCommentsFromTarget(r.Context(), commentID, models.CommentTarget, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comments, http.StatusOK)
}
func (ch *implCommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	idPath := r.PathValue("comment_id")
	commentID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id format", http.StatusBadRequest)
		return
	}
	if err := ch.commentService.DeleteCommentByID(r.Context(), commentID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (ch *implCommentHandler) GetCommentsFromUser(w http.ResponseWriter, r *http.Request) {
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
	comments, err := ch.commentService.GetCommentsFromUser(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, comments, http.StatusOK)
}
