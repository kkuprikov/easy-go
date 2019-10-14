package main

import (
    "github.com/gorilla/websocket"
    "github.com/gomodule/redigo/redis"
    "fmt"
    "strings"
    "log"
    "net/http"
)

type message struct {
  id string
}

const (
  WS_PORT = ":1234"
)

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

func main() {

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
      upgrader.CheckOrigin = func(r *http.Request) bool {
        return true
      }
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
          log.Println(err)
          return
      }
      go Reader(conn, r.URL.Path)
    })

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

        fmt.Println("Received back from client: " + msg.id)
        
        StoreData(path, pool.Get(), msg)
    }
    conn.Close()
}

func StoreData(path string, conn redis.Conn, msg message) {
  parts := strings.Split(path, "/")
  endpoint := parts[1]
  jwt := parts[2]
  fmt.Println(endpoint)
  fmt.Println(jwt)
}