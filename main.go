package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/kyai/redis-cui/cui"
)

var (
	redisHost string
	redisPort string
	redisAuth string
)

func init() {
	flag.StringVar(&redisHost, "h", "127.0.0.1", "Server hostname")
	flag.StringVar(&redisPort, "p", "6379", "Server port")
	flag.StringVar(&redisAuth, "a", "", "Password to use when connecting to the server")
}

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "--help":
			flag.Usage()
			os.Exit(0)
		case "--version":
			fmt.Println(cui.VERSION)
			os.Exit(0)
		default:
			flag.Parse()
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisAuth,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	cui.InitRedisClient(client)

	cui.New()
}
