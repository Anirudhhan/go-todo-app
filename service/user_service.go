package service

import "todo-app/database/dbHelper"

func IsUserExists(email string) (bool, error) {
	return dbHelper.IsUserExists(email)
}

func RegisterUser(name string, email string, passwordHash string) (string, error) {
	return dbHelper.RegisterUser(name, email, passwordHash)
}

func CreateUserSession(userId string) (string, error) {
	return dbHelper.CreateUserSession(userId)
}

func GetUserIDByEmail(email string) (string, string, error) {
	return dbHelper.GetUserIDByEmail(email)
}
