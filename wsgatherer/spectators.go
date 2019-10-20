package wsgatherer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	period          = 3
	periodsToExpire = 3
)

type params struct {
	ID    string
	Count string
}

type message struct {
	Message       string
	Params        params
	BroadcastedAt string
}

func (s *Server) spectatorHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ws, err := wsUpgrade(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 1. save spectator
		// 2. feed back spectators count
		// 3. when spectator leaves, delete
		spectatorProcess(ws, params.ByName("id"), s.Db)
	}
}

func spectatorProcess(ws *websocket.Conn, id string, pool *redis.Pool) {
	saveSpectator(id, pool)
	spectatorFeed(ws, id, pool)
}

func spectatorFeed(ws *websocket.Conn, id string, pool *redis.Pool) {
	ticker := time.NewTicker(period * time.Second)
	defer ticker.Stop()

	for {
		for range ticker.C {
			count, err := spectatorCount(id, false, pool)
			if err != nil {
				fmt.Println("Redis connection error: ", err)
				return
			}

			params := params{
				ID:    id,
				Count: count,
			}
			msg := message{
				Message:       "spectators",
				Params:        params,
				BroadcastedAt: time.Now().String(),
			}

			res, err := json.Marshal(msg)

			if err != nil {
				fmt.Println("JSON marshalling error: ", err)
				return
			}

			if err := ws.WriteMessage(1, res); err != nil {
				deleteSpectator(id, pool)
				return
			}

			err = SendAndClose(pool, "EXPIRE", "{realtime_api}spectators_"+id, periodsToExpire*period)
			if err != nil {
				fmt.Println("Redis connection error: ", err)
			}
		}
	}
}

func saveSpectator(id string, pool *redis.Pool) {
	fmt.Println("Saving spectator to redis")

	_, err := DoAndClose(pool, "INCR", "{realtime_api}spectators_"+id)

	if err != nil {
		fmt.Println("Redis connection error: ", err)
	}
}

func deleteSpectator(id string, pool *redis.Pool) {
	fmt.Println("Spectator left")

	_, err := DoAndClose(pool, "DECR", "{realtime_api}spectators_"+id)

	if err != nil {
		fmt.Println("Redis connection error: ", err)
	}
}

func spectatorCount(id string, withPrefix bool, pool *redis.Pool) (string, error) {
	var prefix string

	if !withPrefix {
		prefix = "{realtime_api}spectators_"
	}

	res, err := redis.String(DoAndClose(pool, "GET", prefix+id))

	if err != nil {
		return "", err
	}

	return res, nil
}

func spectatorsTotal(pool *redis.Pool) ([]byte, error) {
	keys, err := redis.Strings(DoAndClose(pool, "KEYS", "{realtime_api}spectators_*"))

	if err != nil {
		fmt.Println("Redis connection error: ", err)
		return nil, err
	}

	fmt.Println("KEYS", keys)

	var total = make(map[string]string, len(keys))

	for _, key := range keys {
		total[key], err = spectatorCount(key, true, pool)
		if err != nil {
			fmt.Println("Redis connection error: ", err)
			return nil, err
		}
	}

	fmt.Println("TOTAL", total)

	res, err := json.Marshal(total)
	if err != nil {
		fmt.Println("JSON marshaling error: ", err)
		return res, err
	}

	return res, nil
}
