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

// DELETE /conversations/:id/messages/:messageId
func (rt *_router) DeleteMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Messaggio non valido"})
		return
	}
	userIdStr := r.Header.Get("Authorization")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	// Solo il mittente può eliminare il proprio messaggio
	err = rt.db.DeleteMessage(conversationId, messageId, userId)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Non autorizzato o errore eliminazione"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /conversations/:id/messages/:messageId/forward
func (rt *_router) ForwardMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	sourceConvId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Conversazione non valida"})
		return
	}
	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Messaggio non valido"})
		return
	}
	userIdStr := r.Header.Get("Authorization")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	var req struct {
		TargetConversationId int `json:"targetConversationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TargetConversationId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "targetConversationId mancante o non valido"})
		return
	}

	// Recupera il messaggio originale
	original, err := rt.db.GetMessageById(sourceConvId, messageId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Messaggio da inoltrare non trovato"})
		return
	}

	// Inoltra il messaggio usando SendMessage (isForwarded: true)
	messages, err := rt.db.SendMessage(req.TargetConversationId, userId, original.Content, original.MediaType, true)
	if err != nil || len(messages) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Errore inoltro messaggio"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages[len(messages)-1])
}
