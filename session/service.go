package session

import (
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	DBInterface interface {
		//Close() error
	}
)

type Service struct {
	db     redis.Cmdable
	Router *mux.Router
	//email  EmailSender
	randomString func(n int) (string, error)
}

func NewServer(client redis.Cmdable) *Service {
	s := Service{
		client,
		mux.NewRouter().StrictSlash(true),
		randomString,
	}
	s.routes()
	return &s
}
func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
