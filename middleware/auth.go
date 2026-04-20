package middleware

import (
	"fmt"
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.GetHeader("session_id")
		if sessionID == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
			ctx.Abort()
			return
		}

		userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
		if err != nil {
			fmt.Println(err.Error())
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
