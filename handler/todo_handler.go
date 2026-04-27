package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(ctx *gin.Context) {
	var newTodoReq models.CreateTodo
	if err := ctx.ShouldBindJSON(&newTodoReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userID := ctx.GetString("userID")

	if newTodoReq.PendingAt != nil && newTodoReq.PendingAt.Before(time.Now()) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("previous date cannot be inserted"), "previous date cannot be inserted")
		return
	}

	todoID, err := dbHelper.CreateTodo(
		userID,
		newTodoReq,
	)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "new todo created successfully",
		"todo_id": todoID,
	})
}

func UpdateTodo(ctx *gin.Context) {
	var updatedTodoReq models.UpdateTodo

	todoID := ctx.Param("todoID")
	userID := ctx.GetString("userID")

	if err := ctx.ShouldBindJSON(&updatedTodoReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	if updatedTodoReq.PendingAt != nil && updatedTodoReq.PendingAt.Before(time.Now()) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("previous date cannot be inserted"), "previous date cannot be inserted")
		return
	}

	if err := dbHelper.UpdateTodo(todoID, userID, updatedTodoReq); err != nil {
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
		//if err.Error() == "todo not found" {
		//	utils.ErrorResponse(ctx, http.StatusNotFound, err, "todo not found") //todo
		//	return
		//}
		//
		utils.ErrorResponse(ctx, http.StatusNotFound, err, "failed to delete todo")
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
		utils.ErrorResponse(ctx, http.StatusNotFound, err, "failed to fetch todo")
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func GetTodosByUserID(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	status := ctx.Query("status")
	searchValue := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	if status != "" && status != models.Completed && status != models.Pending && status != models.Incomplete {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid status"), "invalid status")
		return
	}

	todos, err := dbHelper.GetTodosByUserID(userID, status, searchValue, page, limit)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"todos": todos,
	})
}
