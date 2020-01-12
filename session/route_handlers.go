package session

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func reportError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	log.Printf("%s", err)
}

func (s *Service) handleCreateSession() http.HandlerFunc {

	type outSession struct {
		SessionId string `json:"session_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var userSessionAttr SessionAttributes
		// Check if userSessionAttr id is provided
		err := json.NewDecoder(r.Body).Decode(&userSessionAttr)
		if err != nil || len(userSessionAttr.UserId) == 0 {
			reportError(w, http.StatusBadRequest, "CreateSession. User is not provided")
			return
		}
		// Create a new session for a provided userSessionAttr id
		var buf bytes.Buffer
		str, err := s.createSession(userSessionAttr, s.sessionExpirationSec)
		if err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to create session in the database. Error: "+err.Error())
			return
		}
		s := outSession{
			str,
		}
		err = json.NewEncoder(&buf).Encode(s)
		if err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to encode session struct to json buffer. Error: "+err.Error())
			return
		}
		// Write response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "json")
		_, err = w.Write(buf.Bytes())
		if err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to write a session struct to a response body. Error: "+err.Error())
		}
	}
}

func (s *Service) handleUpdateAuthInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Service) handleGetSessionAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]
		var user SessionAttributes
		err := s.readSession(sessionId, &user)
		if err != nil {
			reportError(w, http.StatusNotFound,
				"Failed to retrieve session struct from database. Error: "+err.Error())
			return
		}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(user)
		if err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to encode session struct to json buffer. Error: "+err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func (s *Service) handleDropSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]
		err := s.deleteSession(sessionId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
