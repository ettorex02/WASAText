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
