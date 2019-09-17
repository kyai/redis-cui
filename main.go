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
	flag.StringVar(&app.RedisHost, "h", "127.0.0.1", "redis's host")
	flag.StringVar(&app.RedisPort, "p", "6379", "redis's port")
	flag.StringVar(&app.RedisAuth, "a", "", "redis's auth")
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
		default:
			flag.Parse()
		}
	}

	if err := redis.NewPool(
		app.RedisHost+":"+app.RedisPort,
		app.RedisAuth,
		app.RedisPoolMaxIdle,
		app.RedisPoolMaxActive,
		app.RedisPoolIdleTimeout,
	); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	cui.New()
}
