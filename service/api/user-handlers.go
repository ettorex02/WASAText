package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione privata per controllo Authorization
func checkAuthorization(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Non autorizzato"})
		return false
	}
	return true
}

func (rt *_router) GetUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	userId := ps.ByName("userId")
	user, err := rt.db.GetUserById(userId)
	if err != nil || user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Utente non trovato"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PATCH /users/:username/photo per cambiare foto profilo
func (rt *_router) SetMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	userId := ps.ByName("userId")
	var req struct {
		PhotoUrl string `json:"photoUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PhotoUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "URL non valido"})
		return
	}
	if err := rt.db.SetMyPhotoById(userId, req.PhotoUrl); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore aggiornamento immagine"})
		return
	}
	user, _ := rt.db.GetUserById(userId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PATCH /users/:username per cambiare username
func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	userId := ps.ByName("userId")
	var req struct {
		NewName string `json:"newName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.NewName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Username non valido"})
		return
	}
	if err := rt.db.SetMyUserNameById(userId, req.NewName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore aggiornamento username"})
		return
	}
	user, _ := rt.db.GetUserById(userId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GET /users/search?q=...
func (rt *_router) SearchUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	query := r.URL.Query().Get("q")
	if len(query) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Query troppo corta"})
		return
	}
	users, err := rt.db.SearchUsers(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore ricerca utenti"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
