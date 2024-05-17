package routes

import (
	"go-todo-list-api/repository"

	"github.com/gin-gonic/gin"
)

// routes

func Routes() {
	route := gin.Default()

	route.GET("/todos", repository.GetAllTodo)
	route.POST("/todos", repository.CreateTodo)
	route.PUT("/todos/:idTodo", repository.UpdateTodo)
	route.DELETE("/todos/:idTodo", repository.DeleteTodo)

	route.Run("localhost:8080")
}
