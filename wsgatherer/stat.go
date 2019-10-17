package wsgatherer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func queueDict() func(string) string {
	innerMap := map[string]string{
		"heatmap": "heatmap_stats",
		"default": "realtime_stats",
	}

	return func(key string) string {
		return innerMap[key]
	}
}

func StatHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if conn, err := WsUpgrade(w, r); err == nil {
		go statReader(conn, params.ByName("jwt"))
	} else {
		fmt.Println(err)
	}
}

func statReader(ws *websocket.Conn, jwtoken string) {
	var conn = NewRedisConn()()

	for {
		var msg map[string]interface{}

		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Received from client: ", msg["id"])

		data := combineData(jwtoken, msg)
		storeData(conn, data)
	}
	conn.Close()
}

func combineData(jwtoken string, input map[string]interface{}) map[string]interface{} {
	var data map[string]interface{}

	inputJSON := ParseJWT(jwtoken)

	if err := json.Unmarshal(inputJSON, &data); err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)

	for k, v := range data {
		input[k] = v
	}

	return input
}

func storeData(conn redis.Conn, input map[string]interface{}) {
	var queue string

	if event, ok := input["event"].(string); ok {
		queue = queueDict()(event)
	} else {
		fmt.Println("Type assertion failed: event is not a string", event)
	}

	if queue == "" {
		queue = queueDict()("default")
	}

	msg, _ := json.Marshal(input)

	if _, err := conn.Do("LPUSH", queue, string(msg)); err != nil {
		fmt.Println("Could not write to redis: ", err)
	}

	// Read for debug
	if res, err := redis.String(conn.Do("LPOP", queue)); err == nil {
		fmt.Println("Read from redis: ", res)
	} else {
		fmt.Println("Could not read from redis", err)
	}
}
