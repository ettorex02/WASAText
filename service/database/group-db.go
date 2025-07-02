package database

import (
	"fmt"

	"github.com/ettorex02/WASAText/service/structures"
)

// addToGroup: crea un nuovo gruppo con nome e membri (usernames)
func (db *appdbimpl) AddToGroup(name string, usernames []string) (*structures.Group, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO groups (name) VALUES (?)`, name)
	if err != nil {
		return nil, err
	}
	groupID64, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	groupID := int(groupID64)

	var members []structures.User
	for _, username := range usernames {
		var user structures.User
		err := db.c.QueryRow(`SELECT id, username, display_name, profile_picture FROM users WHERE username = ?`, username).
			Scan(&user.ID, &user.Username, &user.DisplayName, &user.ProfilePicture)
		if err != nil {
			return nil, fmt.Errorf("utente %s non trovato", username)
		}
		_, err = tx.Exec(`INSERT INTO group_members (group_id, user_id) VALUES (?, ?)`, groupID, user.ID)
		if err != nil {
			return nil, err
		}
		members = append(members, user)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &structures.Group{
		ID:      groupID,
		Name:    name,
		Photo:   "",
		Members: members,
	}, nil
}

// listGroups: restituisce tutti i gruppi dell'utente autenticato
func (db *appdbimpl) ListGroups(userID int) ([]*structures.Group, error) {
	rows, err := db.c.Query(`
        SELECT g.id, g.name, COALESCE(g.photo, '')
        FROM groups g
        JOIN group_members gm ON g.id = gm.group_id
        WHERE gm.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*structures.Group
	for rows.Next() {
		var group structures.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Photo); err != nil {
			return nil, err
		}
		members, err := db.getGroupMembers(group.ID)
		if err != nil {
			return nil, err
		}
		group.Members = members
		groups = append(groups, &group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

// getGroupByID: restituisce il gruppo completo di membri
func (db *appdbimpl) getGroupByID(groupID int) (*structures.Group, error) {
	var group structures.Group
	err := db.c.QueryRow(`SELECT id, name, COALESCE(photo, '') FROM groups WHERE id = ?`, groupID).
		Scan(&group.ID, &group.Name, &group.Photo)
	if err != nil {
		return nil, err
	}
	members, err := db.getGroupMembers(groupID)
	if err != nil {
		return nil, err
	}
	group.Members = members
	return &group, nil
}

// getGroupMembers: restituisce i membri di un gruppo
func (db *appdbimpl) getGroupMembers(groupID int) ([]structures.User, error) {
	rows, err := db.c.Query(`
        SELECT u.id, u.username, u.display_name, COALESCE(u.profile_picture, '')
        FROM users u
        JOIN group_members gm ON u.id = gm.user_id
        WHERE gm.group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []structures.User
	for rows.Next() {
		var user structures.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.ProfilePicture); err != nil {
			return nil, err
		}
		members = append(members, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}
