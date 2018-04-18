package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"Seckill/pkg/e"
	"Seckill/pkg/util"
	"Seckill/models"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	isExist := models.CheckAuth(username, password)
	if isExist {
		token, err := util.GenerateToken(username, password)
		if err != nil {
			code = e.ERROR_AUTH_TOKEN
		} else {
			data["token"] = token

			code = e.SUCCESS
		}

	} else {
		code = e.ERROR_AUTH
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
