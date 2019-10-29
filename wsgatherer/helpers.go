// Package wsgatherer - this file contains some helper functions
package wsgatherer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	readTimeout = 30 // seconds
)

func readJSON(ws *websocket.Conn) (map[string]string, error) {
	var msg map[string]string

	if err := ws.ReadJSON(&msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func readControl(cancel func(), ws *websocket.Conn) {
	for {
		if _, _, err := ws.NextReader(); err != nil {
			fmt.Println("Client left, canceling context...")
			cancel()

			return
		}
	}
}

func writeControl(ws *websocket.Conn) {
	msg := message{
		Message:       "reconnect",
		BroadcastedAt: time.Now().String(),
	}

	res, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("JSON marshalling error: ", err)
		return
	}

	if err = ws.WriteMessage(1, res); err != nil {
		fmt.Println("Can't write message to websocket: ", err)
		return
	}

	err = ws.SetReadDeadline(time.Now().Add(readTimeout * time.Second))

	if err != nil {
		fmt.Println("Websocket set deadline error: ", err)
		return
	}

	resp, err := readJSON(ws)

	if err != nil {
		switch err.(type) {
		case *websocket.CloseError:
			fmt.Println("Websocket close error: ", err)
			return
		default:
			fmt.Println("Error while reading JSON from client: ", err)
			return
		}
	}

	if resp["message"] == "ok" {
		fmt.Println("Client received and processed reconnect message")
	}
}
