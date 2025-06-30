package database

import "fmt"

// CreateConversation crea una nuova conversazione tra due utenti e restituisce l'id della conversazione
func (db *appdbimpl) CreateConversation(user1, user2 int) (int64, error) {
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	// Controllo esistenza
	var count int
	err := db.c.QueryRow(
		`SELECT COUNT(*) FROM conversations WHERE user1_id = ? AND user2_id = ?`, user1, user2,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, fmt.Errorf("conversazione già esistente")
	}
	// Creazione
	res, err := db.c.Exec(
		`INSERT INTO conversations (user1_id, user2_id) VALUES (?, ?)`, user1, user2,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
