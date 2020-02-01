package session

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
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
	srv := NewServer(client, DefaultSessionExpirationSec)
	srv.randomString = testRandomString

	return mr, srv, nil
}
func TestCreateSession(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := SessionAttributes{
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
	var userInfo SessionAttributes
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
	srv := NewServer(nil, DefaultSessionExpirationSec)
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
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := SessionAttributes{
		"userOne",
		"Netscape on Nokia phone",
		"",
		"",
		"",
		"",
	}
	id, errId := srv.createSession(user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", "/session/"+id, nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GetSession. Invalid HTTP response. Wait %d got %d ", http.StatusOK, w.Code)
	}
	var userInfo SessionAttributes
	err = json.NewDecoder(w.Body).Decode(&userInfo)
	if err != nil {
		t.Errorf("GetSession. Failed to decode GET /session/{id} http response. Error: %s", err.Error())
	}
	if userInfo.UserId != user.UserId || userInfo.DeviceId != user.DeviceId {
		t.Errorf("Corrupted data in the database.")
	}
	mr.Close()
}
func TestReadSessionFailures(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := SessionAttributes{
		"userOne",
		"Netscape on Nokia phone",
		"",
		"",
		"",
		"",
	}
	id, errId := srv.createSession(user, 3600)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	err := srv.readSession("id", &user)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid session id but it didn't.")
	}
	dummy := struct {
		Whatever string `json:"whatever"`
	}{""}
	err = srv.readSession(id, &dummy)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid structure but it didn't.")
	}
	dummyOneValid := struct {
		Whatever string `json:"whatever"`
		UserId   string `json:"user_id"`
	}{"", ""}
	err = srv.readSession(id, &dummyOneValid)
	if err != nil {
		t.Errorf("Expecting readSession not to fail but it did.")
	}
	mr.Close()
}
func TestGetSessionAttributesFailures(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := SessionAttributes{
		"userOne",
		"Netscape on Nokia phone",
		"",
		"",
		"",
		"",
	}
	_, errId := srv.createSession(user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", "/session/"+"id", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("GetSession. Invalid HTTP response. Wait %d got %d ", http.StatusNotFound, w.Code)
	}
	mr.Close()
}
func TestSessionIncompleteUserInfo(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := struct {
		UserId string `json:"user_id"`
	}{
		"userOne",
	}
	id, errId := srv.createSession(user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", "/session/"+id, nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GetSession. Invalid HTTP response. Wait %d got %d ", http.StatusOK, w.Code)
	}
	var userInfo SessionAttributes
	err = json.NewDecoder(w.Body).Decode(&userInfo)
	if err != nil {
		t.Errorf("GetSession. Failed to decode GET /session/{id} http response. Error: %s", err.Error())
	}
	if userInfo.UserId != user.UserId {
		t.Errorf("Failed to read incomplete session attributes from the database.")
	}
	mr.Close()
}
func TestDropSessionDatabaseOnly(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := struct {
		UserId string `json:"user_id"`
	}{
		"userOne",
	}
	id, errId := srv.createSession(user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	err := srv.readSession(id, &user)
	if err != nil {
		t.Error(err)
	}
	err = srv.deleteSession("id")
	if err == nil {
		t.Errorf("Expected an error whe deleteing non existing key but got no error")
	}
	err = srv.deleteSession(id)
	if err != nil {
		t.Error(err)
	}
	err = srv.readSession(id, &user)
	if err == nil {
		t.Errorf("Expected to fail retrieving session attributes but session was retrieved.")
	}
	mr.Close()
}

func TestDropSession(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := struct {
		UserId string `json:"user_id"`
	}{
		"userOne",
	}
	id, errId := srv.createSession(user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("DELETE", "/session/"+"id", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Errorf("Expected an error when deleteing non existing key but got StatusOK")
	}
	req, err = http.NewRequest("DELETE", "/session/"+id, nil)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected no error when deleteing existing key but got %d", w.Code)
	}
	mr.Close()
}

func TestUpdateSessionDB(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	sessionInit := struct {
		UserId string `json:"user_id"`
	}{
		"userOne",
	}
	id, errId := srv.createSession(sessionInit, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	sessionUpd := struct {
		UserId      string `json:"user_id"`
		AccessToken string `json:"access_token"`
	}{
		"userOne",
		"ABC",
	}
	sessionRes := sessionUpd
	sessionRes.AccessToken = ""
	sessionRes.UserId = ""
	err := srv.updateSession(id, sessionUpd, &sessionRes)
	if err != nil {
		t.Errorf("Update session failed %s", err.Error())
	}
	if sessionRes.UserId != sessionUpd.UserId || sessionRes.AccessToken != sessionUpd.AccessToken {
		t.Errorf("Returned session does not match the update %s != %s", sessionUpd, sessionRes)
	}
	if err := srv.updateSession("sessionId", sessionUpd, &sessionRes); err == nil {
		t.Error("Expecting a call fail due to invalid session ID but it didn't")
	}
	mr.Close()
}
