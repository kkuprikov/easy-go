package main

import (
    "golang.org/x/net/websocket"
    "fmt"
    "log"
    "net/http"
)

const WS_PORT = ":1234"

func main() {
    http.Handle("/", websocket.Handler(Echo))

    if err := http.ListenAndServe(WS_PORT, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func Echo(ws *websocket.Conn) {
    var err error

    for {
        var msg string

        if err = websocket.Message.Receive(ws, &msg); err != nil {
            fmt.Println("Can't receive")
            break
        }

        fmt.Println("Received back from client: " + msg)
        StoreData(ws.Request().URL, msg)
    }
}

func RedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}