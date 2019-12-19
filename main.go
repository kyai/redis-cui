package main

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/kyai/redis-cui/cmd"
	"github.com/kyai/redis-cui/cui"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     cmd.RedisHost + ":" + cmd.RedisPort,
		Password: cmd.RedisAuth,
		DB:       cmd.RedisDB,
	})
	if err := client.Ping().Err(); err != nil {
		fmt.Println(err)
		return
	}
	cui.InitRedisClient(client)

	cui.New()
}
