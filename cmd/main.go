// Package main for program start
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kkuprikov/easy-go/wsgatherer"

	"github.com/julienschmidt/httprouter"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &wsgatherer.Server{}

	host, _ := getenvStr("REDIS_HOST")

	port, err := getenvStr("REDIS_PORT")

	if err != nil {
		fmt.Println("Redis port is incorrect, switching to default: ", err)

		port = "6379"
	}

	size, err := getenvInt("REDIS_POOL_SIZE")

	if err != nil {
		size = 10000
	}

	s.Db = wsgatherer.RedisPool(host+":"+port, size)
	s.Router = httprouter.New()
	go s.Start(ctx)

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan
	fmt.Println("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context

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
