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

func SendAndClose(pool *redis.Pool, command string, args ...interface{}) error {
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer conn.Close()

	return conn.Send(command, args...)
}

func DoAndClose(pool *redis.Pool, command string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	defer conn.Close()

	return conn.Do(command, args...)
}
