package dbHelper

import (
	"errors"
	"time"
	"todo-app/database"
	"todo-app/models"

	"github.com/jmoiron/sqlx"
)

func GetAllUsers() ([]models.UserResult, error) {
	query := `SELECT id, name, email, created_at, role, archived_at, suspended_at
				FROM users ORDER BY created_at DESC`

	users := make([]models.UserResult, 0)
	err := database.DB.Select(&users, query)

	return users, err
}

func GetTodos(status string, search string, page, limit int) ([]models.Todo, error) {
	query := `SELECT id,
		       user_id,
		       name,
		       description,
		       pending_at,
		       completed_at,
		       created_at
		FROM todos
		WHERE archived_at IS NULL
		  AND (
		        $1 = ''
		        OR (
		            $1 = 'completed' AND completed_at IS NOT NULL
		        )
		        OR (
		            $1 = 'pending'
		            AND completed_at IS NULL
		            AND (pending_at IS NULL OR pending_at > NOW())
		        )
		        OR (
		            $1 = 'incomplete'
		            AND completed_at IS NULL
		            AND pending_at IS NOT NULL
		            AND pending_at < NOW()
		        )
		  )
		  AND (
		        $2 = ''
		        OR name ILIKE '%' || $2 || '%'
		        OR description ILIKE '%' || $2 || '%'
		  )
		ORDER BY created_at DESC LIMIT $3 OFFSET $4`

	offset := (page - 1) * limit

	todos := make([]models.Todo, 0)
	err := database.DB.Select(&todos, query, status, search, limit, offset)

	return todos, err
}

func UpdateUserSuspension(tx *sqlx.Tx, userID string, suspended bool) error {
	query := `
		UPDATE users
		SET suspended_at =
			CASE
				WHEN $1 = false THEN NULL
				WHEN suspended_at IS NULL THEN NOW()
				ELSE suspended_at
			END
		WHERE id = $2
	`

	res, err := tx.Exec(query, suspended, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func ArchiveUserSessions(tx *sqlx.Tx, userID string) error {
	query := `
		UPDATE user_session
		SET archived_at = $1
		WHERE user_id = $2
		  AND archived_at IS NULL
	`

	_, err := tx.Exec(query, time.Now(), userID)
	return err
}
