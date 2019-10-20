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
	err := conn.Send(command, args...)
	conn.Close()
	return err
}

func DoAndClose(pool *redis.Pool, command string, args ...interface{}) (reply interface{}, err error) {
	conn := pool.Get()
	res, err := conn.Do(command, args...)
	conn.Close()
	return res, err
}
