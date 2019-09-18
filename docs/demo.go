// Just to generate demo data
// Usage:
// go run demo.go

package main

import (
	"encoding/json"
	"time"

	"github.com/kyai/redis-cui/redis"
)

func main() {
	redis.NewPool("127.0.0.1:6379", "", 10, 10, 10)

	conn := redis.Pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")

	// string
	conn.Do("SET", "test:key", "hello world")

	// string/json
	project := struct {
		ID     int       `json:"id"`
		Name   string    `json:"name"`
		Place  string    `json:"place"`
		People int       `json:"people"`
		Time   time.Time `json:"time"`
	}{1, "Hello", "China", 13, time.Now()}
	b, _ := json.Marshal(project)
	conn.Do("SET", "project", string(b))

	// list
	conn.Do("RPUSH", "test:list", 10001, 10002, 10003, 10004, 10005)

	// hash
	conn.Do("HSET", "website", "google", "www.google.com")

	// set
	conn.Do("SADD", "languages", "Java", "Golang", "Python", "PHP")

	// zset
	market := map[string]int{
		"Apple":     11000,
		"Amazon":    9620,
		"Microsoft": 8830,
		"Alphabet":  8390,
		"Facebook":  4600,
		"Alibaba":   4120,
		"Tencent":   3830,
		"Samsung":   2970,
		"Cisco":     2240,
		"Intel":     2220,
	}
	for k, v := range market {
		conn.Do("ZADD", "global:market", v, k)
	}
}
