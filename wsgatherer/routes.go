// Package wsgatherer routes
package wsgatherer

import (
	"log"
	"net/http"
)

// Start method for our server
func (s *Server) Start() {
	s.Router.GET("/", s.testPage())
	s.Router.GET("/info", s.infoPage())
	s.Router.GET("/spectators", s.spectatorsData())

	s.Router.GET("/ws/send_stat/:jwt", s.statHandler())
	s.Router.GET("/ws/subscribe/spectators/:id", s.spectatorHandler())

	if err := http.ListenAndServe(port, s.Router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
