package session

import (
	"net/http"
)

func (s *Server) routes() {
	s.Router.Methods("POST").Name("CreateSession").
		Path("/api/session").HandlerFunc(s.handleCreateSession())
	s.Router.Methods("DELETE").Name("DropSession").
		Path("/api/session/{id}").HandlerFunc(s.handleDropSession())
	s.Router.Methods("GET").Name("GetSessionAttributes").
		Path("/api/session/{id}").HandlerFunc(s.handleGetSessionAttributes())
}

func (s *Server) handleCreateSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (s *Server) handleGetSessionAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"result\":\"OK\"}"))
	}
}

func (s *Server) handleDropSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
