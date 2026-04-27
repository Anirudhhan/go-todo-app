package dbHelper

import (
	"todo-app/database"
	"todo-app/models"
)

func IsUserExist(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users
				WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL AND suspended_at IS NULL;`

	var userExist bool
	err := database.DB.Get(&userExist, query, email)
	return userExist, err
}

func GetUserIDByActiveSession(sessionID string) (string, error) {
	query := `SELECT user_id 
		FROM user_session 
		WHERE id = $1 AND archived_at IS NULL`

	var userID string
	err := database.DB.Get(&userID, query, sessionID)
	return userID, err
}

func RegisterUser(name string, email string, passwordHash string) (string, error) {
	query := `INSERT INTO users(name, email, password) VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id`

	var userID string
	err := database.DB.Get(&userID, query, name, email, passwordHash)

	return userID, err
}

func CreateUserSession(userID string) (string, error) {
	query := `INSERT INTO user_session(user_id) VALUES ($1) RETURNING id`

	var sessionID string
	err := database.DB.Get(&sessionID, query, userID)
	return sessionID, err
}

func GetLoginDetailsByEmail(email string) (models.LoginUserDetails, error) {
	query := `SELECT id, password, role
			FROM users
			WHERE email = TRIM(LOWER($1))
			AND archived_at IS NULL AND suspended_at IS NULL`

	var userDetails models.LoginUserDetails

	err := database.DB.Get(&userDetails, query, email)
	return userDetails, err
}

func ArchiveUserSession(sessionID string) error {
	query := `UPDATE user_session SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL AND suspended_at IS NULL`

	_, err := database.DB.Exec(query, sessionID)
	return err
}

func GetUserDetailsByActiveSession(sessionID string) (models.GetUserDetailsByActiveSessionResult, error) {
	query := `
		SELECT u.id, u.role
		FROM user_session us
		JOIN users u ON u.id = us.user_id
		WHERE us.id = $1
		  AND us.archived_at IS NULL
		  AND u.archived_at IS NULL
		  AND u.suspended_at IS NULL
	`

	var result models.GetUserDetailsByActiveSessionResult
	err := database.DB.Get(&result, query, sessionID)

	return result, err
}
