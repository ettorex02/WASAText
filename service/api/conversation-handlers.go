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
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	var req createConvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "userId mancante"}); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	// userId del chiamante
	user1, _ := strconv.Atoi(r.Header.Get("Authorization"))
	user2 := req.UserId

	convID, err := rt.db.CreateConversation(user1, user2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"message": err.Error()}); encErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Cambiato da StatusConflict
	json.NewEncoder(w).Encode(map[string]interface{}{"conversationId": convID})
}

// GET /conversations
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
