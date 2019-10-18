package wsgatherer

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	port = ":1234"
)

func wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	return conn, err
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/index.html")
}

func Run() {
	router := httprouter.New()
	router.GET("/", homePage)

	router.GET("/ws/send_stat/:jwt", statHandler)
	router.GET("/ws/subscribe/spectators/:id", spectatorHandler)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
