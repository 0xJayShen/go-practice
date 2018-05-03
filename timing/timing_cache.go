package timing

import (
	"fmt"
	"github.com/robfig/cron"
	"gin-docker-mysql/models"
	"gin-docker-mysql/cache"
)

func CacheHotArticle() {
	var articleList []*models.Article
	models.DB.Model(&models.Article{}).Where("state = ?", 1).Find(&articleList)
	for _, article := range articleList {
		_, err := cache.RedisPool.Get().Do("HMSET", article.ID, "tittle", article.Title, "content", article.Content, "creater_time", article.CreatedOn, "modified_time", article.ModifiedOn)

		if err != nil {
			fmt.Println(err)
		}
		//hvalall, _ := redis.StringMap(c.Do("hgetall", hash_key))
		//articl_title ,err:=  redis.String(RedisPool.Get().Do("HGET",article.ID,"tittle"))
		//fmt.Println(articl_title)
		//if err!=nil{
		//	fmt.Println(err)
	}
}

func TimingCache() {
	c := cron.New()
	spec := "：0 0 1 * * ?" //每天凌晨缓存一下
	c.AddFunc(spec, func() {
		CacheHotArticle()
	})
	c.Start()
	select {}
}
