package cache

import (
	"github.com/garyburd/redigo/redis"

	"time"
	"gin-docker-mysql/pkg/setting"
	"fmt"
)

//type Pool struct {
//	//Dial 是创建链接的方法
//	Dial func() (Conn, error)
//
//	//TestOnBorrow 是一个测试链接可用性的方法
//	TestOnBorrow func(c Conn, t time.Time) error
//
//	// 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
//	MaxIdle int
//
//	// 最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
//	MaxActive int
//
//	//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
//	IdleTimeout time.Duration
//
//	// 当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
//	Wait bool
//
//}
//声明一些全局变量
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
	c := RedisPool.Get()
	defer c.Close()
	// _,err = c.Do("Set", "abc", 100)
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


