package handler

import (
	"errors"
	"net/http"
	"time"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(ctx *gin.Context) {
	var createTodo models.CreateTodo
	if err := ctx.ShouldBindJSON(&createTodo); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userID := ctx.GetString("userID")

	if createTodo.PendingAt != nil && createTodo.PendingAt.Before(time.Now()) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("previous date cannot be inserted"), "previous date cannot be inserted")
		return
	}

	todoID, err := dbHelper.CreateTodo(
		userID,
		createTodo.Name,
		createTodo.Description,
		createTodo.PendingAt,
	)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "new todo created successfully",
		"todoId":  todoID,
	})
}

func UpdateTodo(ctx *gin.Context) {
	var updatedTodo models.UpdateTodo

	todoID := ctx.Param("todoID")
	userID := ctx.GetString("userID")

	if err := ctx.ShouldBindJSON(&updatedTodo); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	if updatedTodo.PendingAt != nil && updatedTodo.PendingAt.Before(time.Now()) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("previous date cannot be inserted"), "previous date cannot be inserted")
		return
	}

	if err := dbHelper.UpdateTodo(todoID, userID, updatedTodo); err != nil {
		//if err.Error() == "todo not found" {
		//	utils.ErrorResponse(ctx, http.StatusNotFound, err, "todo not found")
		//	return
		//}
		utils.ErrorResponse(ctx, http.StatusNotFound, err, "failed to update todo")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodo(ctx *gin.Context) {
	todoID := ctx.Param("todoID")
	userID := ctx.GetString("userID")

	err := dbHelper.DeleteTodo(todoID, userID)
	if err != nil {
		if err.Error() == "todo not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, err, "todo not found") //todo
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to update todo")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "todo deleted successfully",
	})
}

func GetTodoByID(ctx *gin.Context) {
	todoID := ctx.Param("todoID")
	userID := ctx.GetString("userID")

	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to fetch todo")
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func GetTodos(ctx *gin.Context) {
	status := ctx.Query("status")
	userID := ctx.GetString("userID")

	todos := make([]models.Todo, 0)
	if status != "" && status != "completed" && status != "pending" && status != "incomplete" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid status"), "invalid status")
		return
	}

	todos, err := dbHelper.GetTodos(userID, status)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, todos)
}
