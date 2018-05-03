package models

import "github.com/jinzhu/gorm"

type Auth struct {
	gorm.Model

	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	DB.Select("id").Where(Auth{Username : username, Password : password}).First(&auth)
	if auth.ID > 0 {

		return true
	}

	return false
}