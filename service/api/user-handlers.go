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
	username := ps.ByName("username")
	user, err := rt.db.GetUserByUsername(username)
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
	username := ps.ByName("username")
	var req struct {
		PhotoUrl string `json:"photoUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PhotoUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "URL non valido"})
		return
	}
	if err := rt.db.SetMyPhoto(username, req.PhotoUrl); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore aggiornamento immagine"})
		return
	}
	user, _ := rt.db.GetUserByUsername(username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PATCH /users/:username per cambiare username
func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	oldUsername := ps.ByName("username")
	var req struct {
		NewName string `json:"newName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.NewName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Username non valido"})
		return
	}
	if err := rt.db.SetMyUserName(oldUsername, req.NewName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore aggiornamento username"})
		return
	}
	user, _ := rt.db.GetUserByUsername(req.NewName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
