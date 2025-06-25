package database

import (
	"database/sql"
	"fmt"

	"github.com/ettorex02/WASAText/service/structures"
)

// CheckUserExistence controlla se esiste un utente con lo username dato
func (db *appdbimpl) CheckUserExistence(username string) (bool, error) {
	var id int
	err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DoLogin controlla se l'utente esiste, se sì restituisce i dati e "login",
// se no lo crea e restituisce i dati e "register".
// Se si tenta di registrare un utente già esistente, restituisce errore.
func (db *appdbimpl) DoLogin(username, displayName, profilePicture string) (*structures.User, string, error) {
	exists, err := db.CheckUserExistence(username)
	if err != nil {
		return nil, "", err
	}

	if exists {
		// Se displayName o profilePicture sono forniti, significa che si sta tentando una doppia registrazione
		if displayName != "" || profilePicture != "" {
			return nil, "", fmt.Errorf("registrazione già effettuata")
		}
		// Recupera i dati dell'utente
		var user structures.User
		err := db.c.QueryRow(
			`SELECT id, username, display_name, profile_picture FROM users WHERE username = ?`, username,
		).Scan(&user.ID, &user.Username, &user.DisplayName, &user.ProfilePicture)
		if err != nil {
			return nil, "", err
		}
		return &user, "login", nil
	}

	// Se manca displayName o profilePicture, non si può registrare
	if displayName == "" || profilePicture == "" {
		return nil, "", fmt.Errorf("per la registrazione servono displayName e profilePicture")
	}

	// Crea nuovo utente
	res, err := db.c.Exec(
		`INSERT INTO users (username, display_name, profile_picture) VALUES (?, ?, ?)`,
		username, displayName, profilePicture,
	)
	if err != nil {
		return nil, "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, "", err
	}
	return &structures.User{
		ID:             int(id),
		Username:       username,
		DisplayName:    displayName,
		ProfilePicture: profilePicture,
	}, "register", nil
}
