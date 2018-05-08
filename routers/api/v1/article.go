package v1

import (
	"net/http"

	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"

	"gin-docker-mysql/models"
	"gin-docker-mysql/pkg/e"
	"gin-docker-mysql/pkg/setting"
	"gin-docker-mysql/pkg/util"
	"fmt"
)

// @Summary 获取单个文章
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{"id":3,"created_on":1516937037,"modified_on":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	var data map[string]interface{}
	data = make(map[string]interface{})

	if models.ExistArticleByID(id) {
		data["article"] = models.GetArticle(id)
		data["comments"] = models.GetComments(id)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_EXIST_ARTICLE
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @Summary 获取多个文章
// @Produce  json
// @Param tag_id query int false "TagID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":[{"id":3,"created_on":1516937037,"modified_on":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1}],"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

	}

	var tagId int = -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
	}

	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

	}

	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

	}

	code := e.SUCCESS

	data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @Summary 新增文章
// @Produce  json
// @Param tag_id query int true "TagID"
// @Param title query string true "Title"
// @Param desc query string true "Desc"
// @Param content query string true "Content"
// @Param created_by query string true "CreatedBy"
// @Param state query int true "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	code := e.INVALID_PARAMS

	if models.ExistTagByID(tagId) {
		data := make(map[string]interface{})
		data["tag_id"] = tagId
		data["title"] = title
		data["desc"] = desc
		data["content"] = content
		data["created_by"] = createdBy
		data["state"] = state

		//models.AddArticle(data)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_EXIST_TAG
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// @Summary 修改文章
// @Produce  json
// @Param id param int true "ID"
// @Param tag_id query string false "TagID"
// @Param title query string false "Title"
// @Param desc query string false "Desc"
// @Param content query string false "Content"
// @Param modified_by query string true "ModifiedBy"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	modifiedBy := c.PostForm("modified_by")

	//var state int = -1
	//if arg := c.PostForm("state"); arg != "" {
	//	state = com.StrTo(arg).MustInt()
	//	maps["state"] = state
	//}

	code := e.INVALID_PARAMS

	if models.ExistArticleByID(id) {
		if models.ExistTagByID(tagId) {
			fmt.Println("cunzai ")
			data := make(map[string]interface{})
			if tagId > 0 {
				data["tag_id"] = tagId
			}
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}

			data["modified_by"] = modifiedBy

			models.EditArticle(id, data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		code = e.ERROR_NOT_EXIST_ARTICLE
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除文章
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS

	if models.ExistArticleByID(id) {
		models.DeleteArticle(id)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_EXIST_ARTICLE
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
