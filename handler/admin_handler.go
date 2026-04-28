package handler

import (
	"errors"
	"net/http"
	"strconv"
	"todo-app/database"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetAllUsersAdmin(ctx *gin.Context) {
	searchValue := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, err := dbHelper.GetAllUsers(searchValue, limit, page)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": len(users),
		"users": users,
	})
}

func GetTodosAdmin(ctx *gin.Context) {

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

	todos, err := dbHelper.GetTodos(status, searchValue, page, limit)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": len(todos),
		"todos": todos,
	})
}

func UpdateUserSuspensionAdmin(ctx *gin.Context) {
	userID := ctx.Param("userID")
	var updateSuspensionRequest struct {
		Suspended bool `json:"suspended"`
	}

	if err := ctx.ShouldBindJSON(&updateSuspensionRequest); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbHelper.UpdateUserSuspension(tx, userID, updateSuspensionRequest.Suspended)
		if err != nil {
			return err
		}

		if updateSuspensionRequest.Suspended {
			err = dbHelper.ArchiveUserSessions(tx, userID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		if txErr.Error() == "user not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, txErr, txErr.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, txErr, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user suspension updated successfully",
	})
}

//func CreateTodoAdmin(ctx *gin.Context) {
//	var newTodoReq models.CreateTodo
//	if err := ctx.ShouldBindJSON(&newTodoReq); err != nil {
//		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
//		return
//	}
//
//	userID := ctx.Param("userID")
//
//	if newTodoReq.PendingAt != nil && newTodoReq.PendingAt.Before(time.Now()) {
//		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("previous date cannot be inserted"), "previous date cannot be inserted")
//		return
//	}
//
//	todoID, err := dbHelper.CreateTodo(userID, newTodoReq)
//	if err != nil {
//		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
//		return
//	}
//
//	ctx.JSON(http.StatusCreated, gin.H{
//		"message": "new todo created successfully",
//		"todo_id": todoID,
//	})
//}
