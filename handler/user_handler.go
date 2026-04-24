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
	var registerUserReq models.RegisterUser

	if err := ctx.ShouldBindJSON(&registerUserReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userExist, err := dbHelper.IsUserExist(registerUserReq.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	if userExist {
		utils.ErrorResponse(ctx, http.StatusConflict, errors.New("user with this email already exists"), "user with this email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(registerUserReq.Password)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	userID, err := dbHelper.RegisterUser(registerUserReq.Name, registerUserReq.Email, hashedPassword)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user_id": userID,
	})
}

func LoginUser(ctx *gin.Context) {
	var loginUserReq models.LoginUser

	if err := ctx.ShouldBindJSON(&loginUserReq); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	userDetails, err := dbHelper.GetLoginDetailsByEmail(loginUserReq.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err, "invalid credentials")
		return
	}

	if err := utils.CheckPasswordHash(loginUserReq.Password, userDetails.HashPassword); err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err, "invalid credentials")
		return
	}

	sessionID, err := dbHelper.CreateUserSession(userDetails.UserID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
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
