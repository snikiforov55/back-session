package session

import (
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
		str, err := s.createSession(userSessionAttr.UserId, userSessionAttr, s.sessionExpirationSec)
		if err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to create session in the database. Error: "+err.Error())
			return
		}
		session := outSession{
			str,
		}
		if js, err := json.Marshal(session); err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to encode session struct to json buffer. Error: "+err.Error())
			return
		} else {
			header := w.Header()
			header.Set("Content-Type", "application/json")
			header.Set("Cache-Control", "no-cache, no-store")
			if _, err := w.Write(js); err != nil {
				reportError(w, http.StatusInternalServerError,
					"Failed write a payload to the response. Error: "+err.Error())
				return
			}
		}
	}
}

func (s *Service) handleUpdateSessionAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//sessionId := vars["id"]
	}
}

func (s *Service) handleGetSessionAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]
		var session SessionAttributes
		err := s.readSession(sessionId, &session)
		if err != nil {
			reportError(w, http.StatusNotFound,
				"Failed to retrieve session struct from database. Error: "+err.Error())
			return
		}
		if js, err := json.Marshal(session); err != nil {
			reportError(w, http.StatusInternalServerError,
				"Failed to encode session struct to json buffer. Error: "+err.Error())
			return
		} else {
			header := w.Header()
			header.Set("Content-Type", "application/json")
			header.Set("Cache-Control", "no-cache, no-store")
			if _, err := w.Write(js); err != nil {
				reportError(w, http.StatusInternalServerError,
					"Failed write a payload to the response. Error: "+err.Error())
				return
			}
		}
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
