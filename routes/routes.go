package routes

import (
	"net/http"
	"todo-app/handler"
	"todo-app/middleware"

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
		v1.PUT("/logout", handler.Logout) //TODO: patch, put or delete?
	}
	{
		todo := v1.Group("/todo")
		todo.Use(middleware.AuthMiddleware())
		{
			todo.POST("/", handler.CreateTodo)
			todo.GET("/", handler.GetAllTodos)
			todo.GET("/:todoID", handler.GetTodoByID)
			todo.PUT("/:todoID", handler.UpdateTodo)
			todo.DELETE("/:todoID", handler.DeleteTodo)
		}
	}

	return routes
}
