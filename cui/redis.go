package cui

import (
	"strconv"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func InitRedisClient(c *redis.Client) {
	client = c
}

func getDatabases() (dbs int, err error) {
	res, err := client.Do("config", "get", "databases").Result()
	if err != nil {
		return
	}
	val, ok := res.([]interface{})
	if !ok {
		return 0, nil
	}
	dbs = 16
	if val[0] == "databases" {
		dbs, _ = strconv.Atoi(val[1].(string))
	}
	return
}
