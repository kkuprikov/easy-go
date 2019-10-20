package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kkuprikov/easy-go/wsgatherer"
)

func main() {
	s := &wsgatherer.Server{}

	host, _ := getenvStr("REDIS_HOST")

	port, err := getenvStr("REDIS_PORT")

	if err != nil {
		port = "6379"
	}

	size, err := getenvInt("REDIS_POOL_SIZE")

	if err != nil {
		size = 10000
	}

	s.Db = wsgatherer.RedisPool(host+":"+port, size)
	s.Router = httprouter.New()
	s.Routes()
}

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getenvInt(key string) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}
