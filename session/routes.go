package session

import (
	"net/http"
)

func (s *Server) routes() {
	s.Router.Methods("POST").Name("InitSession").
		Path("/api/session").HandlerFunc(s.handleInitSession())
	s.Router.Methods("DELETE").Name("DropSession").
		Path("/api/session/{id}").HandlerFunc(s.handleDropSession())
	s.Router.Methods("GET").Name("IsUserAuth").
		Path("/api/session").HandlerFunc(s.handleAllSessions())
}

func (s *Server) handleInitSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (s *Server) handleAllSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"result\":\"OK\"}"))
	}
}

func (s *Server) handleDropSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) handleIsUserAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) handleUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
