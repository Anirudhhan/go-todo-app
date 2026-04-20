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
	sessionID := ctx.GetHeader("session_id")
	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	var createTodo models.CreateTodo
	if err := ctx.ShouldBindJSON(&createTodo); err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
	if err != nil {
		fmt.Println(err.Error())
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

	sessionID := ctx.GetHeader("session_id")
	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	if err := ctx.ShouldBindJSON(&updatedTodo); err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	todoValid, err := dbHelper.IsTodoValid(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if !todoValid {
		utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
		return
	}

	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusNotFound, "internal server error")
		return
	}

	if updatedTodo.Name != "" {
		todo.Name = updatedTodo.Name
	}

	if updatedTodo.Description != "" {
		todo.Description = updatedTodo.Description
	}

	if updatedTodo.PendingAt != nil {
		if updatedTodo.PendingAt.Before(time.Now()) {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid deadline")
			return
		}
		todo.PendingAt = updatedTodo.PendingAt
	}

	if updatedTodo.CompletedAt != nil {
		todo.CompletedAt = updatedTodo.CompletedAt
	}

	if err := dbHelper.UpdateTodo(todoID, userID, todo); err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodo(ctx *gin.Context) {
	todoID := ctx.Param("todoID")
	sessionID := ctx.GetHeader("session_id")

	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	todoValid, err := dbHelper.IsTodoValid(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if !todoValid {
		utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
		return
	}

	err = dbHelper.DeleteTodo(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "todo delete successfully",
	})
}

func GetTodoByID(ctx *gin.Context) {
	todoID := ctx.Param("todoID")
	sessionID := ctx.GetHeader("session_id")

	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	todoValid, err := dbHelper.IsTodoValid(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if !todoValid {
		utils.ErrorResponse(ctx, http.StatusNotFound, "todo not found")
		return
	}

	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func GetAllTodos(ctx *gin.Context) {
	sessionID := ctx.GetHeader("session_id")
	status := ctx.Query("status")

	if sessionID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	userID, err := dbHelper.GetUserIDFromActiveSession(sessionID)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid session")
		return
	}

	var todos []models.Todo
	if status != "completed" && status != "pending" && status != "incomplete" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid status")
		return
	}

	todos, err = dbHelper.GetAllTodos(userID, status)

	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, todos)
}
