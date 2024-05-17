package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// function connect db
func ConnnectDB() *gorm.DB {
	errorENV := godotenv.Load()
	if errorENV != nil {
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect postgres database")
	}

	return db
}

// function Disconnect DB to stopping your connection to postgres database

func disconnectDB(db *gorm.DB) {
	dbSQl, err := db.DB()
	if err != nil {
		panic("Failed to kill connection from database")
	}
	dbSQl.Close()
}

var db *gorm.DB = ConnnectDB()

// todo struct

type Todo struct {
	gorm.Model
	Name        string
	Description string
}

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
	var todos []Todo

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
	todo := Todo{}
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
	todo := Todo{}

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
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went Wrong",
		})
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
	todo := Todo{}

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

// routes

func Routes() {
	route := gin.Default()

	route.GET("/todos", GetAllTodo)
	route.POST("/todos", CreateTodo)
	route.PUT("/todos/:idTodo", UpdateTodo)
	route.DELETE("/todos/:idTodo", DeleteTodo)

	route.Run("localhost:8080")
}

func main() {

	defer disconnectDB(db)

	Routes()
}
