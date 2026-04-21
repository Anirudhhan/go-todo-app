package dbHelper

import (
	"errors"
	"time"
	"todo-app/database"
	"todo-app/models"
)

func CreateTodo(userID string, name string, description string, pendingAt *time.Time) (string, error) {
	query := `INSERT INTO todos(user_id, name, description, pending_at)
				VALUES ($1, TRIM($2), $3, $4)
				RETURNING id`

	var todoID string
	err := database.DB.Get(&todoID, query, userID, name, description, pendingAt)

	return todoID, err
}

func GetTodoByID(todoID string, userID string) (models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todos
				WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var todo models.Todo

	err := database.DB.Get(&todo, query, todoID, userID)
	return todo, err
}

func UpdateTodo(todoID string, userID string, updatedTodo models.UpdateTodo) error {
	query := `UPDATE todos
		SET name = COALESCE(TRIM($1), name), 
		    description = COALESCE(TRIM($2), description), 
		    pending_at = COALESCE($3, pending_at), 
		    completed_at = COALESCE($4, completed_at)  --todo
		WHERE id = $5 AND user_id = $6 AND archived_at IS NULL`

	res, err := database.DB.Exec(query, updatedTodo.Name, updatedTodo.Description, updatedTodo.PendingAt, updatedTodo.CompletedAt, todoID, userID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func DeleteTodo(todoID string, userID string) error {
	query := `UPDATE todos SET archived_at = NOW() WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	res, err := database.DB.Exec(query, todoID, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func GetTodos(userID string, status string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
				FROM todos
				WHERE user_id = $1 AND archived_at IS NULL`

	switch status {
	case "completed":
		query += " AND completed_at IS NOT NULL"

	case "pending":
		query += ` AND completed_at IS NULL AND (pending_at IS NULL OR pending_at > NOW())`

	case "incomplete":
		query += ` AND completed_at IS NULL AND pending_at IS NOT NULL AND pending_at < NOW()`
	}

	query += " ORDER BY created_at DESC"

	todos := make([]models.Todo, 0)
	err := database.DB.Select(&todos, query, userID)

	return todos, err
}
