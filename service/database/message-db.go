package database

import (
	"fmt"

	"github.com/ettorex02/WASAText/service/globaltime"
	"github.com/ettorex02/WASAText/service/structures"
)

// SendMessage inserisce un nuovo messaggio e restituisce la lista aggiornata dei messaggi della conversazione
func (db *appdbimpl) SendMessage(conversationId, senderId int, content, mediaType string, isForwarded bool) ([]*structures.Message, error) {
	// Controlla che la conversazione esista
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, conversationId).Scan(&exists)
	if err != nil || exists == 0 {
		return nil, fmt.Errorf("conversazione non trovata")
	}
	// Controlla che il mittente sia membro della conversazione
	var count int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM conversation_members WHERE conversation_id = ? AND user_id = ?`, conversationId, senderId).Scan(&count)
	if err != nil || count == 0 {
		return nil, fmt.Errorf("utente non autorizzato a inviare messaggi in questa conversazione")
	}

	status := "received"
	_, err = db.c.Exec(
		`INSERT INTO messages (conversation_id, sender_id, content, is_forwarded, media_type, status, timestamp)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		conversationId, senderId, content, isForwarded, mediaType, status, globaltime.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}
	return db.GetMessages(conversationId)
}

// GetMessages restituisce tutti i messaggi di una conversazione, ordinati dal più vecchio al più recente
func (db *appdbimpl) GetMessages(conversationId int) ([]*structures.Message, error) {
	// Controlla che la conversazione esista
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, conversationId).Scan(&exists)
	if err != nil || exists == 0 {
		return nil, fmt.Errorf("conversazione non trovata")
	}

	rows, err := db.c.Query(
		`SELECT m.id, m.conversation_id, m.sender_id, m.content, m.is_forwarded, m.media_type, m.status, m.timestamp,
                u.username, u.display_name, u.profile_picture
         FROM messages m
         JOIN users u ON m.sender_id = u.id
         WHERE m.conversation_id = ?
         ORDER BY m.timestamp ASC`, conversationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var messages []*structures.Message
	for rows.Next() {
		var msg structures.Message
		var sender structures.User
		if err := rows.Scan(
			&msg.ID, &msg.ConversationID, &sender.ID, &msg.Content, &msg.IsForwarded, &msg.MediaType, &msg.Status, &msg.Timestamp,
			&sender.Username, &sender.DisplayName, &sender.ProfilePicture,
		); err != nil {
			return nil, err
		}
		msg.Sender = sender
		messages = append(messages, &msg)
	}
	return messages, nil
}

// Setta a "received" tutti i messaggi ricevuti dall'utente in una conversazione che sono ancora "sent"
func (db *appdbimpl) SetMessagesReceived(conversationId int, userId int) error {
	_, err := db.c.Exec(`
        UPDATE messages
        SET status = CASE WHEN status = 'sent' THEN 'received' ELSE status END
        WHERE conversation_id = ? AND sender_id != ?`, conversationId, userId)
	return err
}

// Setta a "read" tutti i messaggi ricevuti dall'utente in una conversazione che sono "received"
func (db *appdbimpl) SetMessagesRead(conversationId int, userId int) error {
	_, err := db.c.Exec(`
        UPDATE messages
        SET status = 'read'
        WHERE conversation_id = ? AND sender_id != ? AND status != 'read'`, conversationId, userId)
	return err
}

// DeleteMessage rimuove un messaggio da una conversazione se l'utente è il mittente
func (db *appdbimpl) DeleteMessage(conversationId, messageId, userId int) error {
	res, err := db.c.Exec(
		`DELETE FROM messages WHERE id = ? AND conversation_id = ? AND sender_id = ?`,
		messageId, conversationId, userId,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return fmt.Errorf("not authorized or message not found")
	}
	return nil
}

// GetMessageById restituisce un messaggio specifico di una conversazione
func (db *appdbimpl) GetMessageById(conversationId, messageId int) (*structures.Message, error) {
	rows, err := db.c.Query(
		`SELECT m.id, m.conversation_id, m.sender_id, m.content, m.is_forwarded, m.media_type, m.status, m.timestamp,
                u.username, u.display_name, u.profile_picture
         FROM messages m
         JOIN users u ON m.sender_id = u.id
         WHERE m.conversation_id = ? AND m.id = ?`, conversationId, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if rows.Next() {
		var msg structures.Message
		var sender structures.User
		if err := rows.Scan(
			&msg.ID, &msg.ConversationID, &sender.ID, &msg.Content, &msg.IsForwarded, &msg.MediaType, &msg.Status, &msg.Timestamp,
			&sender.Username, &sender.DisplayName, &sender.ProfilePicture,
		); err != nil {
			return nil, err
		}
		msg.Sender = sender
		return &msg, nil
	}
	return nil, fmt.Errorf("messaggio non trovato")
}
