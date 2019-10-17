package wsgatherer

import "github.com/gomodule/redigo/redis"

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func NewRedisConn() func() redis.Conn {
	pool := newPool()

	return func() redis.Conn {
		c := pool.Get()
		return c
	}
}
