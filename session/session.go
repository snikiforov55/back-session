package session

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

type (
	DBInterface interface {
		//Close() error
	}
)

type Server struct {
	db     redis.Cmdable
	Router *mux.Router
	//email  EmailSender
}

func NewServer(client redis.Cmdable) *Server {
	s := Server{
		client,
		mux.NewRouter().StrictSlash(true),
	}
	s.routes()
	return &s
}
