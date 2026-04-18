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

func DeleteTodo(todoID string, userID string) error {
	query := `UPDATE todo SET archived_at = NOW() WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	_, err := database.DB.Exec(query, todoID, userID)
	return err
}

func IsTodoValid(todoID string, userID string) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM todo 
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL)`

	var exists bool
	err := database.DB.Get(&exists, query, todoID, userID)

	return exists, err
}

func GetAllTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todo
				WHERE user_id = $1 AND archived_at IS NULL`

	var todos []models.Todo
	err := database.DB.Select(&todos, query, userID)

	return todos, err
}

func GetAllCompletedTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todo
				WHERE user_id = $1 AND archived_at IS NULL AND completed_at IS NOT NULL`

	var todos []models.Todo
	err := database.DB.Select(&todos, query, userID)

	return todos, err
}

func GetAllPendingTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todo
				WHERE user_id = $1 AND archived_at IS NULL AND completed_at IS NULL AND pending_at > NOW()`

	var todos []models.Todo
	err := database.DB.Select(&todos, query, userID)

	return todos, err
}

func GetAllInCompleteTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todo
				WHERE user_id = $1 AND archived_at IS NULL AND completed_at IS NULL AND pending_at < NOW()`

	var todos []models.Todo
	err := database.DB.Select(&todos, query, userID)

	return todos, err
}
