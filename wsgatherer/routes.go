// Package wsgatherer routes
package wsgatherer

import (
	"context"
	"log"
	"net/http"
)

// Start method for our server
func (s *Server) Start(ctx context.Context) {
	s.Router.GET("/", s.testPage())
	s.Router.GET("/info", s.infoPage())
	s.Router.GET("/spectators", s.spectatorsData())

	s.Router.GET("/ws/send_stat/:jwt", s.statHandler(ctx))
	s.Router.GET("/ws/subscribe/spectators/:id", s.spectatorHandler(ctx))

	if err := http.ListenAndServe(port, s.Router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
