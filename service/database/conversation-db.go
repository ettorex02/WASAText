package database

import (
	"fmt"
	"sort"

	"github.com/ettorex02/WASAText/service/structures"
)

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

func (db *appdbimpl) GetUserConversations(userId int) ([]*structures.ConversationPreview, error) {
	// Prendi tutte le conversazioni dove l'utente è coinvolto
	rows, err := db.c.Query(`
        SELECT id, user1_id, user2_id
        FROM conversations
        WHERE user1_id = ? OR user2_id = ?`, userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previews []*structures.ConversationPreview
	for rows.Next() {
		var id, user1, user2 int
		if err := rows.Scan(&id, &user1, &user2); err != nil {
			return nil, err
		}
		otherId := user2
		if user2 == userId {
			otherId = user1
		}
		// Prendi info dell'altro utente
		var username, profilePicture string
		err = db.c.QueryRow(`SELECT username, profile_picture FROM users WHERE id = ?`, otherId).Scan(&username, &profilePicture)
		if err != nil {
			return nil, err
		}
		// Prendi ultimo messaggio (se esiste)
		var lastMsg, lastTime string
		_ = db.c.QueryRow(`
            SELECT content, timestamp FROM messages
            WHERE conversation_id = ?
            ORDER BY timestamp DESC LIMIT 1
        `, id).Scan(&lastMsg, &lastTime)

		previews = append(previews, &structures.ConversationPreview{
			ID:              id,
			OtherUserID:     otherId,
			Username:        username,
			ProfilePicture:  profilePicture,
			LastMessage:     lastMsg,
			LastMessageTime: lastTime,
		})
	}

	// Ordina in ordine cronologico inverso (ultimo messaggio più recente)
	sort.Slice(previews, func(i, j int) bool {
		return previews[i].LastMessageTime > previews[j].LastMessageTime
	})

	return previews, nil
}
