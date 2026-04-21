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

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		auth.PUT("/logout", handler.Logout)

		{
			todo := auth.Group("/todo")
			{
				todo.POST("/", handler.CreateTodo)
				todo.GET("/", handler.GetTodos)
				todo.GET("/:todoID", handler.GetTodoByID)
				todo.PUT("/:todoID", handler.UpdateTodo)
				todo.DELETE("/:todoID", handler.DeleteTodo)
			}
		}
	}

	return routes
}
