package main

import (
	"go-todo-list-api/config"
	"go-todo-list-api/routes"

	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.ConnnectDB()
)

func main() {

	defer config.DisconnectDB(db)

	routes.Routes()
}
