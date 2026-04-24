package models

import "time"

// status
const (
	Completed  = "completed"
	Pending    = "pending"
	Incomplete = "incomplete"
)

type Todo struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	PendingAt   *time.Time `db:"pending_at" json:"pending_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	//ArchivedAt  *time.Time `db:"archived_at" json:"archived_at"`
}

type CreateTodo struct {
	Name        string     `db:"name" json:"name" binding:"required"`
	Description string     `db:"description" json:"description" binding:"required"`
	PendingAt   *time.Time `db:"pending_at" json:"pending_at" binding:"required"`
}

type UpdateTodo struct {
	Name        *string    `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	PendingAt   *time.Time `db:"pending_at" json:"pending_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
}
