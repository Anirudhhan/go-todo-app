package dbHelper

import (
	"time"
	"todo-app/database"
)

func IsUserExists(email string) (bool, error) {
	query := `SELECT count(*) > 0 FROM users WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL;`

	var exists bool
	err := database.DB.Get(&exists, query, email)
	return exists, err
}

func RegisterUser(name string, email string, passwordHash string) (string, error) {
	query := `INSERT INTO users(name, email, password) VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id`

	var userID string
	err := database.DB.Get(&userID, query, name, email, passwordHash)

	return userID, err
}

func CreateUserSession(userId string) (string, error) {
	query := `INSERT INTO user_session(user_id) VALUES ($1) RETURNING id`

	var sessionId string
	err := database.DB.Get(&sessionId, query, userId)
	return sessionId, err
}

func GetUserIDByEmail(email string) (string, string, error) {
	query := `
		SELECT id, password
		FROM users
		WHERE email = TRIM(LOWER($1))
		AND archived_at IS NULL
	`

	var userID string
	var passwordHash string

	err := database.DB.QueryRow(query, email).Scan(&userID, &passwordHash)
	return userID, passwordHash, err
}

func ArchiveUserSession(sessionId string) error {
	query := `UPDATE user_session SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL`

	_, err := database.DB.Exec(query, sessionId)
	return err
}

func IsSessionActive(sessionID string) (bool, error) {
	query := `SELECT archived_at FROM user_session WHERE id = $1`

	var archivedAt *time.Time
	err := database.DB.Get(&archivedAt, query, sessionID)

	if err != nil {
		return false, err
	}
	return archivedAt == nil, err
}
