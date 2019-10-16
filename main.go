package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	WS_PORT = ":1234"
)

var prefix_to_queue = map[string]string{
	"heatmap": "heatmap_stats",
	"default": "realtime_stats",
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pool = NewPool()

func WsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go Reader(conn, params.ByName("jwt"))
}

func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	router := httprouter.New()
	router.GET("/", HomePage)

	router.GET("/ws/send_stat/:jwt", WsHandler)

	if err := http.ListenAndServe(WS_PORT, router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func Reader(ws *websocket.Conn, jwtoken string) {
	var conn = pool.Get()
	for {
		var msg map[string]interface{}

		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Received from client: ", msg["id"])

		data := CombineData(jwtoken, msg)
		StoreData(pool.Get(), data)
	}
	conn.Close()
}

func CombineData(jwtoken string, input map[string]interface{}) map[string]interface{} {
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

func StoreData(conn redis.Conn, input map[string]interface{}) {
	var queue_name string

	if event, ok := input["event"].(string); ok {
		queue_name = prefix_to_queue[event]
	} else {
		fmt.Println("Type assertion failed: event is not a string", event)
	}

	if queue_name == "" {
		queue_name = prefix_to_queue["default"]
	}

	msg, _ := json.Marshal(input)
	conn.Do("LPUSH", queue_name, msg)

	res, err := redis.String(conn.Do("LPOP", queue_name))
	if err == nil {
		fmt.Println(res)
	}
}
