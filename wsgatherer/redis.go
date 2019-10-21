// Package wsgatherer - this file contains redis communication functions
package wsgatherer

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// RedisPool provides redis pool
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

func sendAndClose(pool *redis.Pool, command string, args ...interface{}) error {
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}

	defer Check(conn.Close)

	return conn.Send(command, args...)
}

func doAndClose(pool *redis.Pool, command string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}

	defer Check(conn.Close)

	return conn.Do(command, args...)
}

// Check helps with checking errors on deferred function calls
func Check(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error: ", err)
	}
}
