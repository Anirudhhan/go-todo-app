package middleware

import (
	"errors"
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("missing token"), "unauthorized")
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid token")
			ctx.Abort()
			return
		}

		sessionID, ok := claims["sid"].(string)
		if !ok {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid claims"), "invalid token")
			ctx.Abort()
			return
		}

		role, _ := claims["role"].(string)
		userID, err := dbHelper.GetUserIDByActiveSession(sessionID)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid session")
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("sessionID", sessionID)
		ctx.Set("role", role)

		ctx.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetString("role") != models.AdminRole {
			utils.ErrorResponse(ctx, http.StatusForbidden, errors.New("forbidden"), "forbidden")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
