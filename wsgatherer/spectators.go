// Package wsgatherer - this files provides API for saving spectators data
package wsgatherer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kkuprikov/easy-go/jcontext"

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

func (s *Server) spectatorHandler(ctx context.Context, wg *sync.WaitGroup) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wg.Add(1)
		defer wg.Done()
		ws, err := wsUpgrade(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 1. save spectator
		// 2. feed back spectators count
		// 3. when spectator leaves, delete

		joinCtx, cancel := jcontext.Join(ctx, r.Context())

		defer cancel()

		go readControl(cancel, ws)
		spectatorProcess(joinCtx, ws, params.ByName("id"), s.Db)
	}
}

func readControl(cancel func(), ws *websocket.Conn) {
	for {
		if _, _, err := ws.NextReader(); err != nil {
			fmt.Println("Client left, canceling context...")
			cancel()

			return
		}
	}
}

func spectatorProcess(ctx context.Context, ws *websocket.Conn, id string, pool *redis.Pool) {
	saveSpectator(id, pool)
	spectatorFeed(ctx, ws, id, pool)
}

func spectatorFeed(ctx context.Context, ws *websocket.Conn, id string, pool *redis.Pool) {
	ticker := time.NewTicker(period * time.Second)
	defer ticker.Stop()
	defer deleteSpectator(id, pool)

	for {
		select {
		//delete spectator on ctx.Done() or reqCtx.Done()
		case <-ctx.Done():
			fmt.Println("ctx.Done() in spectatorFeed")

			writeControl(ws)
			Check(ws.Close)

			return
		case <-ticker.C:
			count, err := spectatorCount(id, false, pool)
			if err != nil {
				fmt.Println("Redis connection error: ", err)

				return
			}

			param := params{
				ID:    id,
				Count: count,
			}
			msg := message{
				Message:       "spectators",
				Params:        param,
				BroadcastedAt: time.Now().String(),
			}

			res, err := json.Marshal(msg)

			if err != nil {
				fmt.Println("JSON marshalling error: ", err)
				return
			}

			if err = ws.WriteMessage(1, res); err != nil {
				fmt.Println("Can't write message to websocket: ", err)
				return
			}

			err = sendAndClose(pool, "EXPIRE", "{realtime_api}spectators_"+id, periodsToExpire*period)
			if err != nil {
				fmt.Println("Redis connection error: ", err)
				return
			}
		}
	}
}

func saveSpectator(id string, pool *redis.Pool) {
	fmt.Println("Saving spectator to redis")

	_, err := doAndClose(pool, "INCR", "{realtime_api}spectators_"+id)

	if err != nil {
		fmt.Println("Redis connection error: ", err)
	}
}

func deleteSpectator(id string, pool *redis.Pool) {
	fmt.Println("Spectator left")

	_, err := doAndClose(pool, "DECR", "{realtime_api}spectators_"+id)

	if err != nil {
		fmt.Println("Redis connection error: ", err)
	}
}

func spectatorCount(id string, withPrefix bool, pool *redis.Pool) (string, error) {
	var prefix string

	if !withPrefix {
		prefix = "{realtime_api}spectators_"
	}

	res, err := redis.String(doAndClose(pool, "GET", prefix+id))

	if err != nil {
		return "", err
	}

	return res, nil
}

func spectatorsTotal(pool *redis.Pool) ([]byte, error) {
	keys, err := redis.Strings(doAndClose(pool, "KEYS", "{realtime_api}spectators_*"))

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
