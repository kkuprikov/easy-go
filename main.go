package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

type message struct {
	ID    string
	Event string
}

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

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

var pool = newPool()

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go Reader(conn, r.URL.Path)
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", Handler)

	if err := http.ListenAndServe(WS_PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func Reader(ws *websocket.Conn, path string) {
	var conn = pool.Get()
	for {
		var msg message

		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Received back from client: " + msg.ID)

		StoreData(path, pool.Get(), msg)
	}
	conn.Close()
}

func StoreData(path string, conn redis.Conn, input message) {
	parts := strings.Split(path, "/")
	endpoint := parts[1]
	jwt := parts[2]

	queue_name := prefix_to_queue[input.Event]
	if queue_name == "" {
		queue_name = prefix_to_queue["default"]
	}

	msg, _ := json.Marshal(input)
	conn.Do("LPUSH", queue_name, msg)

	fmt.Println(endpoint)
	fmt.Println(jwt)
	fmt.Println(conn.Do("LLEN", queue_name))
}
