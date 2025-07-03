package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type createConvRequest struct {
	UserId int `json:"userId"`
}

// POST /conversations
func (rt *_router) CreateConversationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	var req createConvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "userId mancante"})
		return
	}
	// userId del chiamante
	user1, _ := strconv.Atoi(r.Header.Get("Authorization"))
	user2 := req.UserId

	convID, err := rt.db.CreateConversation(user1, user2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Cambiato da StatusConflict
	json.NewEncoder(w).Encode(map[string]interface{}{"conversationId": convID})
}

// GET /conversations
func (rt *_router) GetUserConversationsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	userId, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Utente non autorizzato"})
		return
	}
	convs, err := rt.db.GetUserConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore recupero conversazioni"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}
