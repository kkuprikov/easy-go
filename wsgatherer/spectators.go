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
	period            = 3
	periods_to_expire = 3
)

type params struct {
	Id    string
	Count string
}

type message struct {
	Message       string
	Params        params
	BroadcastedAt string
}

func spectatorHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	conn, err := wsUpgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 1. save spectator
	// 2. feed back spectators count
	// 3. when spectator leaves, delete
	go spectatorProcess(conn, params.ByName("id"))
}

func spectatorProcess(ws *websocket.Conn, id string) {
	saveSpectator(id)
	spectatorFeed(ws, id)
}

func spectatorFeed(ws *websocket.Conn, id string) {
	ticker := time.NewTicker(period * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			count, err := spectatorCount(id, false)
			if err != nil {
				fmt.Println("Redis connection error: ", err)
				return
			}

			var params params
			var msg message

			params.Id = id
			params.Count = count

			msg.Message = "spectators"
			msg.Params = params
			msg.BroadcastedAt = time.Now().String()

			res, err := json.Marshal(msg)

			if err != nil {
				fmt.Println("JSON marshalling error: ", err)
				return
			}

			if err := ws.WriteMessage(1, res); err != nil {
				deleteSpectator(id)
				done <- true
			} else {
				conn := redisConn()()
				conn.Do("EXPIRE", "{realtime_api}spectators_"+id, periods_to_expire*period)
				conn.Close()
			}
		}
	}
}

func saveSpectator(id string) {
	fmt.Println("Saving spectator to redis")
	conn := redisConn()()
	conn.Do("INCR", "{realtime_api}spectators_"+id)
	conn.Close()
}

func deleteSpectator(id string) {
	fmt.Println("Spectator left")
	conn := redisConn()()
	conn.Do("DECR", "{realtime_api}spectators_"+id)
	conn.Close()
}

func flushSpectators() {
	// TODO implement or make redis keys expirable
	return
}

func spectatorCount(id string, with_prefix bool) (string, error) {
	var prefix string
	if with_prefix == false {
		prefix = "{realtime_api}spectators_"
	}
	conn := redisConn()()
	var err error

	res, err := redis.String(conn.Do("GET", prefix+id))
	conn.Close()
	return res, err

}

func spectatorsTotal() ([]byte, error) {
	var total = map[string]string{}
	var err error

	conn := redisConn()()

	keys, err := redis.Strings(conn.Do("KEYS", "{realtime_api}spectators_*"))
	conn.Close()

	if err != nil {
		fmt.Println("Redis connection error", err)
		return nil, err
	}
	fmt.Println("KEYS", keys)

	for _, key := range keys {
		total[key], err = spectatorCount(key, true)
		if err != nil {
			fmt.Println("Redis connection error", err)
			return nil, err
		}
	}
	fmt.Println("TOTAL", total)

	return json.Marshal(total)
}
