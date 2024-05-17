package routes

import (
	"go-todo-list-api/controllers"

	"github.com/gin-gonic/gin"
)

// routes

func Routes() {
	route := gin.Default()

	route.GET("/todos", controllers.GetAllTodo)
	route.POST("/todos", controllers.CreateTodo)
	route.PUT("/todos/:idTodo", controllers.UpdateTodo)
	route.DELETE("/todos/:idTodo", controllers.DeleteTodo)

	route.Run("localhost:8080")
}
