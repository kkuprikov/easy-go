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

func testPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/index.html")
}

func infoPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./static/info.html")
}

func spectatorsData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := spectatorsTotal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(data)
}

func Run() {
	flushSpectators()
	router := httprouter.New()
	router.GET("/", testPage)
	router.GET("/info", infoPage)
	router.GET("/spectators", spectatorsData)

	router.GET("/ws/send_stat/:jwt", statHandler)
	router.GET("/ws/subscribe/spectators/:id", spectatorHandler)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
