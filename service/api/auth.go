package api

import (
	"encoding/json"
	"net/http"
)

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
