// Package wsgatherer - this file contains some helper functions
package wsgatherer

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func readJSON(ws *websocket.Conn) (map[string]string, error) {
	var msg map[string]string

	if err := ws.ReadJSON(&msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func writeControl(ws *websocket.Conn) {
	err := ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""), time.Time{})
	if err != nil {
		fmt.Println("Can't write message to websocket: ", err)
		return
	}

	if _, err := readJSON(ws); err != nil {
		switch err.(type) {
		case *websocket.CloseError:
			fmt.Println("Websocket close error: ", err)
		default:
			fmt.Println("Error while reading JSON from client: ", err)
		}
	}
}
