package session

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

type Server struct {
	db     *redis.Client
	Router *mux.Router
	//email  EmailSender
}

func NewServer() *Server {
	s := Server{
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		mux.NewRouter().StrictSlash(true),
	}
	s.routes()
	return &s
}
