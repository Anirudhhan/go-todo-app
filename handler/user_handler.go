package handler

import (
	"net/http"
	"todo-app/models"
	"todo-app/service"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	var newUser models.RegisterUser

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	exists, err := service.IsUserExists(newUser.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if exists {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "user with this email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	id, err := service.RegisterUser(newUser.Name, newUser.Email, hashedPassword)
	if err != nil {
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
		utils.ErrorResponse(ctx, http.StatusBadRequest, bindErr.Error())
		return
	}

	userId, hashPassword, err := service.GetUserIDByEmail(loginUser.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid credentials")
		return
	}

	if hashErr := utils.CheckPasswordHash(loginUser.Password, hashPassword); hashErr != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid credentials")
		return
	}

	sessionId, SessionErr := service.CreateUserSession(userId)
	if SessionErr != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessionId": sessionId,
	})
}

func Logout(ctx *gin.Context) {
	sessionId := ctx.Param("sessionId")

	active, err := service.IsSessionActive(sessionId)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid session")
		return
	}

	if !active {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "user already logged out")
		return
	}

	err = service.ArchiveUserSession(sessionId)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid session")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out successfully",
	})
}
