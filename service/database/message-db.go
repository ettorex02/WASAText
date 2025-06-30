package database

import (
	"github.com/ettorex02/WASAText/service/globaltime"
	"github.com/ettorex02/WASAText/service/structures"
)

// SendMessage inserisce un nuovo messaggio e restituisce la lista aggiornata dei messaggi della conversazione
func (db *appdbimpl) SendMessage(conversationId, senderId int, content, mediaType string, isForwarded bool) ([]*structures.Message, error) {
	_, err := db.c.Exec(
		`INSERT INTO messages (conversation_id, sender_id, content, is_forwarded, media_type, status, timestamp)
         VALUES (?, ?, ?, ?, ?, 'sent', ?)`,
		conversationId, senderId, content, isForwarded, mediaType, globaltime.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}
	return db.GetMessages(conversationId)
}

// GetMessages restituisce tutti i messaggi di una conversazione, ordinati dal più vecchio al più recente
func (db *appdbimpl) GetMessages(conversationId int) ([]*structures.Message, error) {
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
