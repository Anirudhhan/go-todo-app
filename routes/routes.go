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
		todo := v1.Group("/todo")
		{
			todo.POST("/", handler.CreateTodo)
			todo.PUT("/:todoID", handler.UpdateTodo)
		}
	}

	return routes
}
