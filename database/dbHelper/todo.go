package dbHelper

import (
	"time"
	"todo-app/database"
)

func CreateTodo(userID string, name string, description string, pendingAt *time.Time) (string, error) {
	query := `INSERT INTO todo(user_id, name, description, pending_at)
				VALUES ($1, TRIM($2), $3, $4)
				RETURNING id`

	var todoID string
	err := database.DB.Get(&todoID, query, userID, name, description, pendingAt)

	return todoID, err
}
