package utils

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hashed), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func ErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"error": message,
	})
}
