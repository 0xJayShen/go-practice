package timing

import (
	"fmt"
	"github.com/robfig/cron"
	"gin-docker-mysql/models"
	"gin-docker-mysql/cache"

	"gin-docker-mysql/pkg/logging"
	"time"
)

func CacheHotArticle() {
	conn := cache.RedisPool.Get()
	defer  conn.Close()
	var articleList []*models.Article
	models.DB.Model(&models.Article{}).Where("state = ?", 1).Find(&articleList)
	for _, article := range articleList {
		_, err := conn.Do("HMSET", article.ID, "tittle", article.Title, "content", article.Content, "creater_time", article.CreatedOn, "modified_time", article.ModifiedOn)

		if err != nil {
			fmt.Println(err)
		}
		//hvalall, _ := redis.StringMap(c.Do("hgetall", hash_key))
		//articl_title ,err:=  redis.String(RedisPool.Get().Do("HGET",article.ID,"tittle"))
		//fmt.Println(articl_title)
		//if err!=nil{
		//	fmt.Println(err)
	}
	logging.Info("每日文章缓存成功 %v",time.Now().Format("2006-01-02 15:04:05"))
}

func TimingCache() {
	c := cron.New()
	spec := "：0 0 1 * * ?" //每天凌晨缓存一下
	//spec := "*/1 * * * * *"
	c.AddFunc(spec, func() {
		CacheHotArticle()
	})
	c.Start()
	select {}
}
