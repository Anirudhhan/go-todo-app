package handler

import (
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var createTodo models.CreateTodo

	if bindErr := ctx.ShouldBindJSON(&createTodo); bindErr != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, bindErr.Error())
		return
	}

	valid, err := dbHelper.IsUserIDValid(userId)
	if !valid || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "userid doesn't exist",
		})
		return
	}

	todoId, err := dbHelper.CreateTodo(userId, createTodo.Name, createTodo.Description, createTodo.PendingAt)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "new todo created successfully",
		"todoId":  todoId,
	})
}
