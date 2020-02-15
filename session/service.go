package session

import (
	"github.com/gorilla/mux"
	"github.com/snikiforov55/back-session/session/db"
	"net/http"
)

type Service struct {
	db     db.Database
	Router *mux.Router
	//email  EmailSender
	sessionExpirationSec int
}

func NewServer(client db.Database, defaultSessionExpSec int) (*Service, error) {
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
