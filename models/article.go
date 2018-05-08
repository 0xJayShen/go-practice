package models

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"gin-docker-mysql/cache"
	"gin-docker-mysql/pkg/logging"
	"time"

	"encoding/json"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
	ReadCount  int    `json:"read_count"`
}

func ExistArticleByID(id int) bool {
	var article Article
	DB.Select("id").Where("id = ? ", id, ).First(&article)

	if article.ID > 0 {
		return true
	}
	return false
}

func GetArticleTotal(maps interface{}) (count int) {

	DB.Model(&Article{}).Where(maps).Count(&count)

	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	DB.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	//Preload就是一个预加载器，它会执行两条SQL，分别是SELECT * FROM blog_articles;和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);，那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，并且避免了循环查询
	//gorm的Join

	return
}

func GetArticle(id int) interface{} {
	conn := cache.RedisPool.Get()
	defer  conn.Close()
	is_key_exit, _ := redis.Bool(conn.Do("EXISTS", id))
	fmt.Println(is_key_exit)
	if is_key_exit {

		//article,_ := redis.StringMap(conn.Do("HGETALL", id))
		//
		//ret = make(map[string]interface{})
		//
		//return article
		var imapGet map[string]interface{}
		valueGet, err := redis.Bytes(conn.Do("GET", id))
		if err != nil {
			fmt.Println(err)
		}
		_ = json.Unmarshal(valueGet, &imapGet)
		fmt.Println(imapGet["tag"])
		fmt.Println(imapGet["title"])

		return imapGet

	} else {
		fmt.Println("没调用缓存")
		var article Article
		DB.Where("id = ?", id).First(&article)
		DB.Model(&article).Related(&article.Tag)
		//Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
		//Article有一个结构体成员是Tag嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
		//fmt.Println(reflect.TypeOf(article))
		//jsons,err:=json.Marshal(article)
		//_, err = conn.Do("HMSET", article.ID,"article",jsons )
		//if err != nil {
		//	fmt.Println(err)
		//	logging.Error("取文章时存缓存出错 %v %v", article.ID, time.Now().Format("2006-01-02 15:04:05"))
		//}
		value := map[string]interface{}{"title":article.Title,"tag":article.Tag,"content":article.Content}
		value_json,_ := json.Marshal(value)
		_,err :=  conn.Do("SET",article.ID,value_json)
		if err != nil {
			fmt.Println(err)
			logging.Error("取文章时存缓存出错 %v %v", article.ID, time.Now().Format("2006-01-02 15:04:05"))
		}

		return article
	}

}

func EditArticle(id int, data interface{}) bool {
	var article Article
	is_key_exit, _ := redis.Bool(cache.RedisPool.Get().Do("EXISTS", id))
	if is_key_exit {
		cache.RedisPool.Get().Do("DEL", id)
		//数据库更新
		DB.Model(&Article{}).Where("id = ?", id).Updates(data)

		//重新写入缓存
		DB.Where("id = ? ", id).First(&article)
		DB.Model(&article).Related(&article.Tag)
		_, err := cache.RedisPool.Get().Do("HMSET", article.ID, "tittle", article.Title, "content", article.Content, "creater_time", article.CreatedOn, "modified_time", article.ModifiedOn)
		if err != nil {
			logging.Error("更新文章缓存出错")
		}
	}

	return true
}

func AddArticle(data map[string]interface{}) bool {
	DB.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
		ReadCount: data["read_count"].(int),
	})

	return true
}

func DeleteArticle(id int) bool {
	is_key_exit, _ := redis.Bool(cache.RedisPool.Get().Do("EXISTS", id))
	if is_key_exit{
		cache.RedisPool.Get().Do("DEL", id)
	}
	DB.Where("id = ?", id).Delete(Article{})
	return true
}
