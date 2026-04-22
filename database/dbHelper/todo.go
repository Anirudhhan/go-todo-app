package dbHelper

import (
	"errors"
	"todo-app/database"
	"todo-app/models"
)

func CreateTodo(userID string, newTodo models.CreateTodo) (string, error) {
	query := `INSERT INTO todos(user_id, name, description, pending_at)
				VALUES ($1, $2, $3, $4)
				RETURNING id`

	var todoID string
	err := database.DB.Get(&todoID, query, userID, newTodo.Name, newTodo.Description, newTodo.PendingAt)

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
		    description = COALESCE($2, description), 
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

func GetTodos(userID, status, search string, page, limit int) ([]models.Todo, error) {
	//query := `SELECT id, user_id, name, description, pending_at, completed_at, created_at, archived_at
	//		  FROM todos
	//		  WHERE user_id = $1 AND archived_at IS NULL`
	//
	//args := []any{userID}
	//argPos := 2
	//
	//switch status {
	//case "completed":
	//	query += " AND completed_at IS NOT NULL"
	//case "pending":
	//	query += ` AND completed_at IS NULL AND (pending_at IS NULL OR pending_at > NOW())`
	//case "incomplete":
	//	query += ` AND completed_at IS NULL AND pending_at IS NOT NULL AND pending_at < NOW()`
	//}
	//
	//if search != "" {
	//	query += ` AND (name ILIKE % || $2 || % OR description ILIKE % || $2 || %)`
	//	args = append(args, search)
	//	argPos++
	//}
	//
	//offset := (page - 1) * limit
	//
	//query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	//
	//args = append(args, limit, offset)
	query := `SELECT id,
		       user_id,
		       name,
		       description,
		       pending_at,
		       completed_at,
		       created_at,
		       archived_at
		FROM todos
		WHERE user_id = $1
		  AND archived_at IS NULL
		  AND (
		        $2 = ''
		        OR (
		            $2 = 'completed' AND completed_at IS NOT NULL
		        )
		        OR (
		            $2 = 'pending'
		            AND completed_at IS NULL
		            AND (pending_at IS NULL OR pending_at > NOW())
		        )
		        OR (
		            $2 = 'incomplete'
		            AND completed_at IS NULL
		            AND pending_at IS NOT NULL
		            AND pending_at < NOW()
		        )
		  )
		  AND (
		        $3 = ''
		        OR name ILIKE '%' || $3 || '%'
		        OR description ILIKE '%' || $3 || '%'
		  )
		ORDER BY created_at DESC LIMIT $4 OFFSET $5`

	offset := (page - 1) * limit

	todos := make([]models.Todo, 0)
	err := database.DB.Select(&todos, query, userID, status, search, limit, offset)

	return todos, err
}
