package middleware

import (
	"errors"
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.GetHeader("Authorization")
		if sessionID == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("invalid session"), "invalid session")
			ctx.Abort()
			return
		}

		userID, err := dbHelper.GetUserIDByActiveSession(sessionID)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, err, "invalid session")
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("sessionID", sessionID)
		ctx.Next()
	}
}
