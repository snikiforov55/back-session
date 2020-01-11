package session

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func reportError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	log.Printf("%s", err)
}

func (s *Service) handleCreateSession() http.HandlerFunc {
	type inUserInfo struct {
		UserId             string `json:"user_id"`
		DeviceId           string `json:"device_id,omitempty"`
		AuthenticationCode string `json:"auth_code"`
		AccessToken        string `json:"access_token"`
		RefreshToken       string `json:"refresh_token"`
		UserEmail          string `json:"user_email"`
	}
	type outSession struct {
		SessionId string `json:"session_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var user inUserInfo //{"", ""}
		// Check if user id is provided
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil || len(user.UserId) == 0 {
			reportError(w, http.StatusBadRequest, "CreateSession. User is not provided")
			return
		}
		// Create a new session for a provided user id
		var buf bytes.Buffer
		str, err := s.randomString(47)
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
		w.Write([]byte("{\"result\":\"OK\"}"))
	}
}

func (s *Service) handleDropSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
