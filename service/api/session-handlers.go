package api

import (
	"encoding/json"
	"net/http"

	"github.com/ettorex02/WASAText/service/structures"
	"github.com/julienschmidt/httprouter"
)

// Assicurati che _router abbia il campo db di tipo AppDatabase
// type _router struct { db database.AppDatabase }

func (rt *_router) SessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req structures.SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, `{"message":"Invalid request: name is required and must be 3-16 characters"}`, http.StatusBadRequest)
		return
	}

	user, action, err := rt.db.DoLogin(req.Name, req.DisplayName, req.ProfilePicture)
	if err != nil {
		if err.Error() == "registrazione già effettuata" {
			http.Error(w, `{"message":"Registrazione già effettuata"}`, http.StatusBadRequest)
			return
		}
		if err.Error() == "per la registrazione servono displayName e profilePicture" {
			http.Error(w, `{"message":"Per la registrazione servono displayName e profilePicture"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	var msg string
	if action == "login" {
		msg = "Login effettuato"
	} else if action == "register" {
		msg = "Registrazione effettuata"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string           `json:"message"`
		User    *structures.User `json:"user"`
	}{
		Message: msg,
		User:    user,
	})
}
