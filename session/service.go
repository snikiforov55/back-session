package session

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	db     redis.Cmdable
	Router *mux.Router
	//email  EmailSender
	sessionExpirationSec int
}

func NewServer(client redis.Cmdable, defaultSessionExpSec int) (*Service, error) {
	s := Service{
		client,
		mux.NewRouter().StrictSlash(true),
		defaultSessionExpSec,
	}
	s.routes()
	return &s, nil
}
func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
