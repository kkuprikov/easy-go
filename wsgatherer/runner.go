package wsgatherer

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
	port = ":1234"
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

var pool = NewPool()

func WsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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

func Run() {
	router := httprouter.New()
	router.GET("/", HomePage)

	router.GET("/ws/send_stat/:jwt", WsHandler)

	if err := http.ListenAndServe(port, router); err != nil {
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
	if _, err := redis.String(conn.Do("LPUSH", queue, msg)); err != nil {
		fmt.Println("Could not write to redis")
	}

	// Read for debug
	if res, err := redis.String(conn.Do("LPOP", queue)); err != nil {
		fmt.Println(res)
	} else {
		fmt.Println("Could not read from redis")
	}
}
