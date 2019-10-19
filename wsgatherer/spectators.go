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

func (s *Server) spectatorHandler() httprouter.Handle {
	conn := s.Db.Get()
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ws, err := wsUpgrade(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 1. save spectator
		// 2. feed back spectators count
		// 3. when spectator leaves, delete
		go spectatorProcess(ws, params.ByName("id"), conn)
	}
}

func spectatorProcess(ws *websocket.Conn, id string, conn redis.Conn) {
	saveSpectator(id, conn)
	spectatorFeed(ws, id, conn)
	defer conn.Close()
}

func spectatorFeed(ws *websocket.Conn, id string, conn redis.Conn) {
	ticker := time.NewTicker(period * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count, err := spectatorCount(id, false, conn)
			if err != nil {
				fmt.Println("Redis connection error: ", err)
				return
			}

			params := params{
				Id:    id,
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
				deleteSpectator(id, conn)
				return
			}
			conn.Do("EXPIRE", "{realtime_api}spectators_"+id, periods_to_expire*period)
		}
	}
}

func saveSpectator(id string, conn redis.Conn) {
	fmt.Println("Saving spectator to redis")
	conn.Do("INCR", "{realtime_api}spectators_"+id)
}

func deleteSpectator(id string, conn redis.Conn) {
	fmt.Println("Spectator left")
	conn.Do("DECR", "{realtime_api}spectators_"+id)
}

func flushSpectators() {
	// TODO implement or make redis keys expirable
	return
}

func spectatorCount(id string, with_prefix bool, conn redis.Conn) (string, error) {
	var prefix string
	if with_prefix == false {
		prefix = "{realtime_api}spectators_"
	}

	res, err := redis.String(conn.Do("GET", prefix+id))
	if err != nil {
		return "", err
	}
	return res, nil

}

func spectatorsTotal(conn redis.Conn) ([]byte, error) {
	keys, err := redis.Strings(conn.Do("KEYS", "{realtime_api}spectators_*"))

	if err != nil {
		fmt.Println("Redis connection error", err)
		return nil, err
	}
	fmt.Println("KEYS", keys)

	var total = make(map[string]string, len(keys))

	for _, key := range keys {
		total[key], err = spectatorCount(key, true, conn)
		if err != nil {
			fmt.Println("Redis connection error", err)
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
