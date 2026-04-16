package models

import "time"

type User struct {
	ID         string     `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"password"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at"`
}

type RegisterUser struct {
	Name     string `db:"name" json:"name" binding:"required,min=3"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Password string `db:"password" json:"password" binding:"required,min=6,max=20"`
}

type LoginUser struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}
