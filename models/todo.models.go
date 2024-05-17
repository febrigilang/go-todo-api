package models

import "gorm.io/gorm"

// todo struct

type Todo struct {
	gorm.Model
	Name        string
	Description string
}
