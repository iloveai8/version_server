package main

import (
	`fmt`
	`github.com/gin-gonic/gin`
	`github.com/gomodule/redigo/redis`
	`time`
	`game_slots_vsn/config`
)

var (
	redisPool *redis.Pool
)

func TestPing(){
	conn := redisPool.Get()
	res, _ := conn.Do("ping")
	fmt.Println(res.(string))
	conn.Close()
}

func newRedisPool(addr string, port string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 2,
		MaxActive: 2,
		IdleTimeout: 6 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr + ":" + port) },
	}
}

func HandlerGetVersionUrl(c *gin.Context) {
		conn := redisPool.Get()
		defer conn.Close()
		conn.Do("Select", 1)
		data_vsn, _ := conn.Do("Get", "vsn")
		data_url, _ := conn.Do("Get", "url")
		version_id := string(data_vsn.([]byte))
		version_url := string(data_url.([]byte))

		var msg version_data
		msg.Vsn = version_id
		msg.Url = version_url
		c.JSON(200, msg)
	
}

type version_data struct {
	Vsn string
	Url string
}

func main() {
	r := gin.Default()
	//读取配置
	cfg, _ := tool.ParseConfig("./config/config.json")
	// 建立redis连接池
	redisPool = newRedisPool(cfg.RedisConfig.Addr, cfg.RedisConfig.Port)
	TestPing()

	r.GET("/get_vesion_url", func(c *gin.Context) {
		HandlerGetVersionUrl(c)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}