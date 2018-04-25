package v1

import (
	"github.com/gin-gonic/gin"
	"gin-docker-mysql/pkg/e"
	"net/http"
	"github.com/Unknwon/com"
	"gin-docker-mysql/models"
	"fmt"
	"gin-docker-mysql/pkg/util"
	"gin-docker-mysql/pkg/setting"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var code int

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		maps["state"] = state
		code = e.SUCCESS

	}else {
		maps["state"] = 1
		code = e.SUCCESS

	}


	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	code := e.SUCCESS
	fmt.Println(name, state, createdBy)
	models.AddTag(name, state, createdBy)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	var state int = -1



	code := e.SUCCESS
	if models.ExistTagByID(id) {
		data := make(map[string]interface{})
		data["modified_by"] = modifiedBy
		if name != "" {
			data["name"] = name
		}
		if state != -1 {
			data["state"] = state
		}

		models.EditTag(id, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	fmt.Println(id)

	code := e.SUCCESS
	if models.ExistTagByID(id) {
		fmt.Println()
		models.DeleteTag(id)
	} else {
		code = e.ERROR_NOT_EXIST_TAG
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
//func A(c *gin.Context){
//	maps := make(map[string]interface{})
//	data := make(map[string]interface{})
//
//	arg := c.Query("state")
//	state := com.StrTo(arg).MustInt()
//	fmt.Println(state)
//	maps["state"] = state
//	data["total"] = models.GetTags(maps)
//	fmt.Println(data)
//	c.JSON(http.StatusOK, gin.H{
//
//		"data": data,
//	})
//}