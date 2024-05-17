package repository

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

//update todos

func UpdateTodo(ctx *gin.Context) {
	var data todoRequest

	// defining request paramater to get id
	reqParamId := ctx.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)

	//Binding request body json to request body struct
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// initiate todo
	todo := models.Todo{}

	//query find todo data by todo id from request parameter
	todoByid := db.Where("id = ?", idTodo).First(&todo)
	if todoByid.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Todo not found",
		})
		return
	}

	// update new todo data
	result := db.Save(&todo)
	fmt.Println(result)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went Wrong",
		})
		return
	}

	// Matching result to todo response struct
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description

	//Creating http response
	ctx.JSON(http.StatusCreated, response)
}

// delete todo

func DeleteTodo(ctx *gin.Context) {
	todo := models.Todo{}

	reqParamId := ctx.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)

	delete := db.Where("id = ?", idTodo).Unscoped().Delete(&todo)
	fmt.Println(delete)

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": "200",
		"message":    "Success",
		"data":       idTodo,
	})
}
