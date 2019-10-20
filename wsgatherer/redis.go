package wsgatherer

import "github.com/gomodule/redigo/redis"

func RedisPool(addr string, size int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: size, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
