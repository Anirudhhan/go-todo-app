package server

import (
	"net/http"
	"todo-app/handler"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes() *gin.Engine {
	routes := gin.Default()
	v1 := routes.Group("v1")

	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "server is running",
			})
		})
		v1.POST("/register", handler.RegisterUser)
		v1.GET("/login", handler.LoginUser)
		v1.PATCH("/logout/:sessionId", handler.Logout) //TODO: patch, edit or delete?
	}

	return routes
}
