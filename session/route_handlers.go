package session

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/snikiforov55/back-session/session/db"
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
		str, err := s.db.CreateSession(userSessionAttr.UserId, userSessionAttr, s.sessionExpirationSec)
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
		payload := struct {
			SessionId         string            `json:"session_id"`
			SessionAttributes map[string]string `json:"session_attributes"`
		}{
			"",
			make(map[string]string),
		}
		//var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			reportError(w, http.StatusUnprocessableEntity,
				"Failed to unmarhal a payload from a request. Error: "+err.Error())
			return
		}
		if payload.SessionId == "" {
			reportError(w, http.StatusBadRequest, "Payload does not contain a session id")
			return
		}
		if payload.SessionAttributes == nil || len(payload.SessionAttributes) == 0 {
			reportError(w, http.StatusBadRequest, "Payload does not contain session attributes to set")
			return
		}
		response := payload.SessionAttributes
		if err := s.db.UpdateSession(payload.SessionId, payload.SessionAttributes, &response); err != nil {
			reportError(w, http.StatusBadRequest, "Failed to update database. "+err.Error())
			return
		}
		if js, err := json.Marshal(response); err != nil {
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

func (s *Service) handleGetSessionAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["id"]
		var session SessionAttributes
		err := s.db.ReadSession(sessionId, &session)
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

func (s *Service) handleGetUserSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["user_id"]
		sessions, err := s.db.ReadUserSessions(userId, []string{"device_id", db.CreateDateAttr, db.UpdateDateAttr})
		if err != nil {
			reportError(w, http.StatusNotFound,
				"Failed to retrieve sessions for user"+userId+" from database. Error: "+err.Error())
			return
		}
		if js, err := json.Marshal(sessions); err != nil {
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
		err := s.db.DeleteSession(sessionId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
