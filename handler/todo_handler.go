package handler

import (
	"net/http"
	"time"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(ctx *gin.Context) {
	sessionID := ctx.GetHeader("session_id")
	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	var createTodo models.CreateTodo
	if err := ctx.ShouldBindJSON(&createTodo); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := dbHelper.GetUserIDFromSession(sessionID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	if createTodo.PendingAt != nil && createTodo.PendingAt.Before(time.Now()) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "previous date cannot be inserted")
		return
	}

	todoID, err := dbHelper.CreateTodo(
		userID,
		createTodo.Name,
		createTodo.Description,
		createTodo.PendingAt,
	)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "new todo created successfully",
		"todoId":  todoID,
	})
}
