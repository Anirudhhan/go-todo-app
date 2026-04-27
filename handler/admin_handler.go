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

func GetAllUsers(ctx *gin.Context) {
	users, err := dbHelper.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func GetTodos(ctx *gin.Context) {

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

func UpdateUserSuspension(ctx *gin.Context) {
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
