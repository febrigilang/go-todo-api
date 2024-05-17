package controllers

import (
	"fmt"
	"go-todo-list-api/config"
	"go-todo-list-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var db *gorm.DB = config.ConnnectDB()

// todo struct request body

type todoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// todo struct response

type todoResponse struct {
	todoRequest
	ID uint `json:"id"`
}

// handles func

// get all todo

func GetAllTodo(ctx *gin.Context) {
	var todos []models.Todo

	//query to find todo datas
	err := db.Find(&todos)
	if err.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error Getting data",
		})
		return
	}

	//create http response

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": "200",
		"message":    "Success",
		"data":       todos,
	})
}

func CreateTodo(ctx *gin.Context) {
	var data todoRequest

	//binding request body json to request body struct
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Matching todo models struct with todo request struct
	todo := models.Todo{}
	todo.Name = data.Name
	todo.Description = data.Description

	// query to database
	result := db.Create(&todo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	// Matching result to create response

	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description

	// create http response
	ctx.JSON(http.StatusCreated, response)
}

// Update todo data
func UpdateTodo(context *gin.Context) {
	var data todoRequest

	// Defining request parameter to get todo id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)

	// Binding request body json to request body struct
	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initiate models todo
	todo := models.Todo{}

	// Querying find todo data by todo id from request parameter
	todoById := db.Where("id = ?", idTodo).First(&todo)
	if todoById.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found"})
		return
	}

	// Matching todo request with todo models
	todo.Name = data.Name
	todo.Description = data.Description

	// Update new todo data
	result := db.Save(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
		return
	}

	// Matching result to todo response struct
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description

	// Creating http response
	context.JSON(http.StatusCreated, response)
}

// Delete todo data function
func DeleteTodo(context *gin.Context) {
	// Initiate todo models
	todo := models.Todo{}
	// Getting request parameter id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)

	// Querying find todo data by todo id from request parameter
	todoById := db.Where("id = ?", idTodo).First(&todo)
	if todoById.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Todo not found"})
		return
	}

	// Querying delete todo by id
	delete := db.Where("id = ?", idTodo).Unscoped().Delete(&todo)
	fmt.Println(delete)

	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    idTodo,
	})

}
