package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kkuprikov/easy-go/wsgatherer"
)

func main() {
	// wsgatherer.Run()
	s := new(wsgatherer.Server)
	s.Db = wsgatherer.RedisPool()
	s.Router = httprouter.New()
	wsgatherer.Routes(s)
}
