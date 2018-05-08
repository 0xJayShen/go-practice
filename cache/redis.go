package cache

import (
	"github.com/garyburd/redigo/redis"

	"time"
	"gin-docker-mysql/pkg/setting"
	"fmt"
)


var (
	RedisPool *redis.Pool
)

func init() {

	var (
		redis_address                      string
		max_idle, max_active, idle_timeout int
	)

	sec, err := setting.Cfg.GetSection("redis")
	if err != nil {
		fmt.Println("error")
	}
	redis_address = sec.Key("RedisAddress").String()
	max_idle = sec.Key("RedisMaxIdle").MustInt(8)
	max_active = sec.Key("RedisMaxActive").MustInt(8)
	idle_timeout = sec.Key("RedisIdleTimeout").MustInt(1)

	RedisPool = &redis.Pool{
		MaxIdle:     max_idle,
		MaxActive:   max_active,
		IdleTimeout: time.Duration(idle_timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redis_address)
		},
	}
	conn :=  RedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		fmt.Println("ping 不通")
		return
	}
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//r, err := redis.Int(c.Do("Get", "abc"))
	//fmt.Println(r)
	//if err != nil {
	//	fmt.Println("get abc failed,", err)
	//	return
	//}
	//redisPool.Close()
}

func CloseRedis() {
	defer RedisPool.Close()
}


