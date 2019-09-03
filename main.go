package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kyai/redis-cui/app"
	"github.com/kyai/redis-cui/cui"
	"github.com/kyai/redis-cui/redis"
)

func init() {
	flag.StringVar(&app.RedisHost, "h", "127.0.0.1:6379", "redis's host")
	flag.StringVar(&app.RedisAuth, "p", "", "redis's auth")
	flag.Parse()
}

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "--help":
			flag.Usage()
			os.Exit(0)
		case "--version":
			fmt.Println(app.VERSION)
			os.Exit(0)
		}
	}

	if err := redis.NewPool(app.RedisHost, app.RedisAuth, 10, 10, 10); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	cui.New()
}
