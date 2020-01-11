package session

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v7"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	session = 0
)

func testRandomString(_ int) (string, error) {
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

func setupServer() (*miniredis.Miniredis, *Service, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	srv := NewServer(client)
	srv.randomString = testRandomString

	return mr, srv, nil
}
func TestCreateSession(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := UserInfo{
		"userOne",
		"Netscape on Nokia phone",
		"",
		"",
		"",
		"",
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
	var userInfo UserInfo
	err = srv.readSession(session.SessionId, &userInfo)
	if err != nil {
		t.Errorf("Failed to read back from database. Error: " + err.Error())
	}
	if userInfo.UserId != user.UserId || userInfo.DeviceId != user.DeviceId {
		t.Errorf("Corrupted data in the database.")
	}
	mr.Close()
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
