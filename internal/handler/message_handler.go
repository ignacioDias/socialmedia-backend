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

type MessageHandler interface {
	CreateChat(w http.ResponseWriter, r *http.Request)
	GetChatsFromUser(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
	GetMessagesFromChat(w http.ResponseWriter, r *http.Request)
	CreateMessage(w http.ResponseWriter, r *http.Request)
}
type implMessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageRepo database.MessageRepository, cache *cache.Cache) MessageHandler {
	return &implMessageHandler{messageService: service.NewMessageService(messageRepo, cache)}
}

func (mh *implMessageHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var body map[string]int64
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	receiver := body["user_id"]
	chat := models.NewChat(userID, receiver)
	if err := mh.messageService.CreateChat(r.Context(), chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, chat, http.StatusCreated)
}

func (mh *implMessageHandler) GetChatsFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	chats, err := mh.messageService.GetChatsFromUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, chats, http.StatusOK)

}

func (mh *implMessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("message_id")
	messageID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := mh.messageService.DeleteMessage(r.Context(), messageID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (mh *implMessageHandler) GetMessagesFromChat(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("chat_id")
	chatID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

}

func (mh *implMessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {

}
