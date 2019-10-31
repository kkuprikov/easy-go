// Package wsgatherer routes
package wsgatherer

import (
	"context"
	"log"
	"net/http"
	"sync"
)

// Start method for our server
func (s *Server) Start(ctx context.Context, port string, wg *sync.WaitGroup) {
	s.Router.GET("/ready", s.ready(ctx))
	s.Router.GET("/spectators", s.spectatorsData())

	s.Router.GET("/ws/send_stat/:jwt", s.statHandler(ctx, wg))
	s.Router.GET("/ws/subscribe/spectators/:id", s.spectatorHandler(ctx, wg))
	s.Router.GET("/", s.assets())

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
