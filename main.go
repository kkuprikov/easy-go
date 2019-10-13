package main

import (
    "golang.org/x/net/websocket"
    "github.com/gomodule/redigo/redis"
    "fmt"
    "log"
    "net/http"
    "net/url"
)

const WS_PORT = ":1234"

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


func main() {

    http.Handle("/", websocket.Handler(Reader))

    if err := http.ListenAndServe(WS_PORT, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func Reader(ws *websocket.Conn) {
    var err error
    pool := newPool()

    for {
        var msg string

        if err = websocket.Message.Receive(ws, &msg); err != nil {
            fmt.Println("Can't receive")
            break
        }

        fmt.Println("Received back from client: " + msg)
        // fmt.Println(reflect.TypeOf(ws.Request().URL))
        StoreData(ws.Request().URL, pool.Get(), msg)
    }
}

func StoreData(url *url.URL, c redis.Conn, msg string) {
  c.Do("SET", "foo", "bar")
  res, _ := redis.String(c.Do("GET", "foo"))
  fmt.Println(string(res))
  defer c.Close()
}