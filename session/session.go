package session

import (
	"github.com/gorilla/mux"
)

type Server struct {
	//db     *someDatabase
	Router *mux.Router
	//email  EmailSender
}

func NewServer() *Server {
	s := Server{
		mux.NewRouter().StrictSlash(true),
	}
	s.routes()
	return &s
}
