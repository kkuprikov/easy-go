// Package wsgatherer handles websocket connections and provides API for data storing
package wsgatherer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// Server struct holds redis database and httprouter
type Server struct {
	Db     *redis.Pool
	Router *httprouter.Router
}

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

func (s *Server) assets() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.FileServer(assetFS())
	}
}

func (s *Server) ready(ctx context.Context) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data struct {
			Status string
		}

		if ctx.Err() != nil {
			data.Status = "Server shutdown"
		} else {
			data.Status = "OK"
		}

		resp, err := json.Marshal(data)

		if err != nil {
			fmt.Println("JSON marshalling error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(resp)
		if err != nil {
			fmt.Println("Error writing a response: ", err)
		}
	}
}

func (s *Server) spectatorsData() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data, err := spectatorsTotal(s.Db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			fmt.Println("Error writing a response: ", err)
		}
	}
}
