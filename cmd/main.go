// Package main for program start
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/kkuprikov/easy-go/wsgatherer"

	"github.com/julienschmidt/httprouter"
)

func main() {
	appPort, err := getenvStr("WSGATHERER_PORT")

	if err != nil {
		fmt.Println("Application port not provided: ", err)
		return
	}

	redisHost, err := getenvStr("REDIS_HOST")

	if err != nil {
		redisHost = "redis"
		fmt.Println("Redis host not provided, switching to: ", redisHost)
		fmt.Println("details: ", err)
	}

	redisPort, err := getenvStr("REDIS_PORT")

	if err != nil {
		redisPort = "6379"
		fmt.Println("Redis port not provided, switching to: ", redisPort)
		fmt.Println("details: ", err)
	}

	size, err := getenvInt("REDIS_POOL_SIZE")

	if err != nil {
		size = 10000
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	s := &wsgatherer.Server{}

	s.Db = wsgatherer.RedisPool(redisHost+":"+redisPort, size)
	s.Router = httprouter.New()

	go s.Start(ctx, appPort, &wg)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan
	fmt.Println("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()

	fmt.Println("All workers done, shutting down!")
}

var errEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", errEnvVarEmpty
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
