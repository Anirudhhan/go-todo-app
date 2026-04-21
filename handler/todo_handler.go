package handler

import (
	"fmt"
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
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID := ctx.GetString("userID")

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
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
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
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := dbHelper.UpdateTodo(todoID, userID, updatedTodo); err != nil {
		if err.Error() == "todo not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update todo")
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
			utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update todo")
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
		if err.Error() == "todo not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
			return
		}

		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update todo")
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func GetAllTodos(ctx *gin.Context) {
	status := ctx.Query("status")
	userID := ctx.GetString("userID")

	todos := make([]models.Todo, 0)
	if status != "" && status != "completed" && status != "pending" && status != "incomplete" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid status")
		return
	}

	todos, err := dbHelper.GetAllTodos(userID, status)

	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, todos)
}
