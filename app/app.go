package app

const (
	VERSION = "0.1.0"
)

const (
	RedisPoolMaxIdle     = 10
	RedisPoolMaxActive   = 10
	RedisPoolIdleTimeout = 10
)

var (
	RedisHost string
	RedisPort string
	RedisAuth string
	RedisDB   int
)
