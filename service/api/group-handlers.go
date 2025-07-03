package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// POST /groups (operationId: addToGroup)
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	var req struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
		Photo   string   `json:"photo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || len(req.Members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Dati gruppo non validi"})
		return
	}
	group, err := rt.db.AddToGroup(req.Name, req.Photo, req.Members)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

// GET /groups (operationId: listGroups)
func (rt *_router) listGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	userID, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Utente non autorizzato"})
		return
	}
	groups, err := rt.db.ListGroups(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore recupero gruppi"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// DELETE /groups/:id/members (operationId: leaveGroup)
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID gruppo non valido"})
		return
	}
	userID, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil || userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Utente non autorizzato"})
		return
	}
	if err := rt.db.LeaveGroup(groupID, userID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /groups/:id/name (operationId: setGroupName)
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID gruppo non valido"})
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Nome non valido"})
		return
	}
	if err := rt.db.SetGroupName(groupID, req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /groups/:id/photo (operationId: setGroupPhoto)
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID gruppo non valido"})
		return
	}
	var req struct {
		Photo string `json:"photo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Photo == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Foto non valida"})
		return
	}
	if err := rt.db.SetGroupPhoto(groupID, req.Photo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /groups/:id/members
func (rt *_router) addGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !checkAuthorization(w, r) {
		return
	}
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID gruppo non valido"})
		return
	}
	var req struct {
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Membri non validi"})
		return
	}
	if err := rt.db.AddMembersToGroup(groupID, req.Members); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Errore aggiunta membri"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
