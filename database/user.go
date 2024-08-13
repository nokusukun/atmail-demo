package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	Age         int    `json:"age"`
	Permissions string `json:"permissions" gorm:"default:'PUT,DELETE'"`
}
