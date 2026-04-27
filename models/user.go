package models

import "time"

type UserResult struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	//Password   string     `db:"password" json:"password"`
	Role        string     `db:"role" json:"role"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	ArchivedAt  *time.Time `db:"archived_at" json:"archived_at"`
	SuspendedAt *time.Time `db:"suspended_at" json:"suspended_at"`
}

type RegisterUser struct {
	Name     string `db:"name" json:"name" binding:"required,min=3"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Password string `db:"password" json:"password" binding:"required,min=6,max=20"`
}

type LoginUserDetails struct {
	UserID       string `db:"id"`
	HashPassword string `db:"password"`
	Role         string `db:"role"`
}

type GetUserDetailsByActiveSessionResult struct {
	UserID string `db:"id"`
	Role   string `db:"role"`
}

type LoginUser struct {
	Email    string `db:"email" json:"email" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}
