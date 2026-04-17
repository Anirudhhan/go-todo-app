package routes

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
		v1.POST("/login", handler.LoginUser)
		v1.PATCH("/logout", handler.Logout) //TODO: patch, put or delete?
	}
	{
		v1.POST("/todo", handler.CreateTodo)
	}

	return routes
}
