package dbHelper

import (
	"time"
	"todo-app/database"
	"todo-app/models"
)

func CreateTodo(userID string, name string, description string, pendingAt *time.Time) (string, error) {
	query := `INSERT INTO todo(user_id, name, description, pending_at)
				VALUES ($1, TRIM($2), $3, $4)
				RETURNING id`

	var todoID string
	err := database.DB.Get(&todoID, query, userID, name, description, pendingAt)

	return todoID, err
}

func GetTodoByID(todoID string, userID string) (models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todo
				WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var todo models.Todo

	err := database.DB.Get(&todo, query, todoID, userID)
	return todo, err
}

func UpdateTodo(todoID string, userID string, updatedTodo models.Todo) error {
	query := `UPDATE todo 
		SET name = $1, description = $2, pending_at = $3, completed_at = $4
		WHERE id = $5 AND user_id = $6 AND archived_at IS NULL`

	_, err := database.DB.Exec(query, updatedTodo.Name, updatedTodo.Description, updatedTodo.PendingAt, updatedTodo.CompletedAt, todoID, userID)
	return err
}
