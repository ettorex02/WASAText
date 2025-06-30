package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// POST /conversations/:id/messages
func (rt *_router) SendMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Conversazione non valida"})
		return
	}
	userId, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Utente non autorizzato"})
		return
	}
	var req struct {
		Content     string `json:"content"`
		MediaType   string `json:"mediaType"`
		IsForwarded bool   `json:"isForwarded"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Contenuto mancante"})
		return
	}
	if req.MediaType == "" {
		req.MediaType = "text"
	}
	message, err := rt.db.SendMessage(conversationId, userId, req.Content, req.MediaType, req.IsForwarded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore invio messaggio"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message) // Restituisci solo il messaggio inviato
}

// GET /conversations/:id/messages
func (rt *_router) GetMessagesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Conversazione non valida"})
		return
	}
	userIdStr := r.Header.Get("Authorization")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Aggiorna a "received" i messaggi ricevuti da questo utente
	_ = rt.db.SetMessagesReceived(conversationId, userId)

	messages, err := rt.db.GetMessages(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Errore recupero messaggi"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// PUT /conversations/:id/messages/read
func (rt *_router) SetMessagesReadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Conversazione non valida"})
		return
	}
	userIdStr := r.Header.Get("Authorization")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	err = rt.db.SetMessagesRead(conversationId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Errore aggiornamento messaggi"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
