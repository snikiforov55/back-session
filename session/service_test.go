package session

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	session = 0
)

func testRandomString(n int) (string, error) {
	return "SESSION_01", nil
}

func TestRandomString(t *testing.T) {
	str, err := randomString(37)
	if err != nil {
		t.Error(err)
	}
	if len(str) < 37 {
		t.Errorf("Unexpected random string length. Waiting for >= 47 got %d", len(str))
	}
}
func TestCreateSession(t *testing.T) {
	srv := NewServer(nil)
	srv.randomString = testRandomString
	user := struct {
		UserId   string `json:"user_id"`
		DeviceId string `json:"device_id"`
	}{
		UserId:   "userOne",
		DeviceId: "Netscape on Nokia phone",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	if err != nil {
		t.Error(err)
	} // json.NewEncoder
	req, err := http.NewRequest("POST", "/session", &buf)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("CreateSession. Invalid HTTP response. Wait %d got %d ", http.StatusOK, w.Code)
	}
	session := struct {
		SessionId string `json:"session_id"`
	}{""}
	err = json.NewDecoder(w.Body).Decode(&session)
	if err != nil {
		t.Errorf("CreateSession. Failed to decode a session id from the HTTP response. "+
			"Error: %s", err)
	}
	if session.SessionId != "SESSION_01" {
		t.Errorf("CreateSession. Unexpected sesson id. "+
			"Waiting for \"SESSION_01\" got \"%s\"", session.SessionId)
	}
}

func TestCreateSessionNoUser(t *testing.T) {
	srv := NewServer(nil)
	user := struct {
		UserId string `json:"user_id"`
	}{
		UserId: "",
	}
	var buf bytes.Buffer
	req, err := http.NewRequest("POST", "/session", &buf)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("CreateSession. Invalid HTTP response. Wait %d got %d ", http.StatusBadRequest, w.Code)
	}
	err = json.NewEncoder(&buf).Encode(user)
	if err != nil {
		t.Error(err)
	} // json.NewEncoder
	req, err = http.NewRequest("POST", "/session", &buf)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("CreateSession. Invalid HTTP response. Wait %d got %d ", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateAuthInfo(t *testing.T) {

}

func TestGetSessionAttributes(t *testing.T) {

}

func TestDropSession(t *testing.T) {

}
