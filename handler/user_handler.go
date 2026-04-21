package handler

import (
	"fmt"
	"net/http"
	"todo-app/database/dbHelper"
	"todo-app/models"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var newUser models.RegisterUser

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	exist, err := dbHelper.IsUserExists(newUser.Email)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if exist {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "user with this email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	id, err := dbHelper.RegisterUser(newUser.Name, newUser.Email, hashedPassword)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"id":      id,
	})
}

func LoginUser(ctx *gin.Context) {
	var loginUser models.LoginUser

	if bindErr := ctx.ShouldBindJSON(&loginUser); bindErr != nil {
		fmt.Println(bindErr.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, bindErr.Error())
		return
	}

	userId, hashPassword, err := dbHelper.GetUserIDByEmail(loginUser.Email)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid credentials")
		return
	}

	if hashErr := utils.CheckPasswordHash(loginUser.Password, hashPassword); hashErr != nil {
		fmt.Println(hashErr.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid credentials")
		return
	}

	sessionId, SessionErr := dbHelper.CreateUserSession(userId)
	if SessionErr != nil {
		fmt.Println(SessionErr.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessionId": sessionId,
	})
}

func Logout(ctx *gin.Context) {
	sessionId := ctx.GetHeader("session_id")
	if sessionId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid session")
		return
	}

	active, err := dbHelper.IsSessionActive(sessionId)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid session")
		return
	}

	if !active {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "user already logged out")
		return
	}

	err = dbHelper.ArchiveUserSession(sessionId)
	if err != nil {
		fmt.Println(err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid session")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out successfully",
	})
}
