package session

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"net/http"
)

type SessionAttributes struct {
	UserId             string `json:"user_id"`
	DeviceId           string `json:"device_id"`
	AuthenticationCode string `json:"auth_code"`
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	UserEmail          string `json:"user_email"`
}

type Service struct {
	db     redis.Cmdable
	Router *mux.Router
	//email  EmailSender
	randomString         func(n int) (string, error)
	sessionExpirationSec int
}

func NewServer(client redis.Cmdable) *Service {
	s := Service{
		client,
		mux.NewRouter().StrictSlash(true),
		randomString,
		3600 * 24 * 365,
	}
	s.routes()
	return &s
}
func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
