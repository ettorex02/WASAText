package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ettorex02/WASAText/service/structures"
	"github.com/julienschmidt/httprouter"
)

var (
	users  = make(map[string]*structures.User)
	nextID = 1
	mu     sync.Mutex
)

type SessionRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
}

func (rt *_router) SessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, `{"message":"Invalid request: name is required and must be 3-16 characters"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	user, exists := users[req.Name]
	if !exists {
		// Per la creazione, servono anche displayName e profilePicture
		if req.DisplayName == "" || req.ProfilePicture == "" {
			http.Error(w, `{"message":"User does not exist. To register, provide also displayName and profilePicture."}`, http.StatusBadRequest)
			return
		}
		user = &structures.User{
			ID:             nextID,
			Username:       req.Name,
			DisplayName:    req.DisplayName,
			ProfilePicture: req.ProfilePicture,
		}
		users[req.Name] = user
		nextID++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			Message string           `json:"message"`
			User    *structures.User `json:"user"`
		}{
			Message: "User created successfully",
			User:    user,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string           `json:"message"`
		User    *structures.User `json:"user"`
	}{
		Message: "Login successful",
		User:    user,
	})
}
