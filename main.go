package main

import (
	"flag"
	"fmt"
	"os"
	"redis-cui/app"
	"redis-cui/cui"
	"redis-cui/redis"
)

var (
	host string
	auth string
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:6379", "redis's host")
	flag.StringVar(&auth, "p", "", "redis's auth")
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

	if err := redis.NewPool(host, auth, 10, 10, 10); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	cui.New()
}
