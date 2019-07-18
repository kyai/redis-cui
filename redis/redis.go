package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Pool      *redis.Pool
	String    = redis.String
	Strings   = redis.Strings
	Int       = redis.Int
	StringMap = redis.StringMap
)

func NewPool(server, password string, maxIdle, maxActive, idleTimeout int) error {
	Pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err = c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	_, err := String(Pool.Get().Do("PING"))
	return err
}
