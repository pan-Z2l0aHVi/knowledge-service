package tools

import (
	"fmt"

	"github.com/go-redis/redis"
)

var rds *redis.Client

type Redis struct {
}

func (e *Redis) InitRedis() {
	cfg, err := ParseConfigure("./conf/app.json")
	if err != nil {
		panic(err)
	}
	address := cfg.Redis.Host + ":" + cfg.Redis.Port
	fmt.Println("address", address)
	rds = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Redis.Pwd,
		DB:       0,
	})
}

func (e *Redis) GetRDS() *redis.Client {
	return rds
}
