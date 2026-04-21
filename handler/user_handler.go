package handler

import (
	"errors"
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var registerUserRes models.RegisterUser

	if err := ctx.ShouldBindJSON(&registerUserRes); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userExist, err := dbHelper.IsUserExists(registerUserRes.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if userExist {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New("user with this email already exist"), "user with this email already exist")
		return
	}

	hashedPassword, err := utils.HashPassword(registerUserRes.Password)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	id, err := dbHelper.RegisterUser(registerUserRes.Name, registerUserRes.Email, hashedPassword)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"id":      id,
	})
}

func LoginUser(ctx *gin.Context) {
	var loginUserRes models.LoginUser

	if bindErr := ctx.ShouldBindJSON(&loginUserRes); bindErr != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, bindErr, bindErr.Error())
		return
	}

	userDetails, err := dbHelper.GetUserIDAndHashedPassByEmail(loginUserRes.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, "invalid credentials")
		return
	}

	if hashErr := utils.CheckPasswordHash(loginUserRes.Password, userDetails.HashPassword); hashErr != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, hashErr, "invalid credentials")
		return
	}

	sessionID, sessionErr := dbHelper.CreateUserSession(userDetails.UserID)
	if sessionErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, sessionErr, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessionId": sessionID,
	})
}

func Logout(ctx *gin.Context) {
	sessionID := ctx.GetString("sessionID")

	err := dbHelper.ArchiveUserSession(sessionID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "failed to logout user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out successfully",
	})
}
