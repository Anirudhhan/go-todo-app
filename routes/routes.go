package routes

import (
	"net/http"
	"todo-app/handler"
	"todo-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "server is running",
			})
		})
		v1.POST("/register", handler.RegisterUser)
		v1.POST("/login", handler.LoginUser)
		v1.GET("/refresh", handler.RefreshToken)

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		auth.PUT("/logout", handler.Logout)

		{
			todo := auth.Group("/todo")
			{
				todo.POST("/", handler.CreateTodo)
				todo.GET("/", handler.GetTodosByUserID)
				todo.GET("/:todoID", handler.GetTodoByID)
				todo.PUT("/:todoID", handler.UpdateTodo)
				todo.DELETE("/:todoID", handler.DeleteTodo)
			}

			admin := auth.Group("/admin")
			admin.Use(middleware.AdminAuthMiddleware())
			{
				admin.GET("/users", handler.GetAllUsers)
				admin.GET("/todos", handler.GetTodos)
				admin.POST("/user/:userID", handler.UpdateUserSuspension)
			}
		}
	}

	return router
}
