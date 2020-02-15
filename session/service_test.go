package session

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/snikiforov55/back-session/session/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRandomString(_ int) (string, error) {
	return "SESSION_01", nil
}
func setupServer() (*miniredis.Miniredis, *Service, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}
	client := db.NewRedisClient(&db.RedisConfig{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}, testRandomString)
	srv, _ := NewServer(client, DefaultSessionExpirationSec)
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
	req, err := http.NewRequest("POST", Api()+"/session", &buf)
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
	err = srv.db.ReadSession(session.SessionId, &userInfo)
	if err != nil {
		t.Errorf("Failed to read back from database. Error: " + err.Error())
	}
	if userInfo.UserId != user.UserId || userInfo.DeviceId != user.DeviceId {
		t.Errorf("Corrupted data in the database.")
	}
	mr.Close()
}
func TestCreateSessionNoUser(t *testing.T) {
	srv, _ := NewServer(nil, DefaultSessionExpirationSec)
	user := struct {
		UserId string `json:"user_id"`
	}{
		UserId: "",
	}
	var buf bytes.Buffer
	req, err := http.NewRequest("POST", Api()+"/session", &buf)
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
	req, err = http.NewRequest("POST", Api()+"/session", &buf)
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
	id, errId := srv.db.CreateSession(user.UserId, user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", Api()+"/session/"+id, nil)
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
	id, errId := srv.db.CreateSession(user.UserId, user, 3600)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	err := srv.db.ReadSession("id", &user)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid session id but it didn't.")
	}
	dummy := struct {
		Whatever string `json:"whatever"`
	}{""}
	err = srv.db.ReadSession(id, &dummy)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid structure but it didn't.")
	}
	dummyOneValid := struct {
		Whatever string `json:"whatever"`
		UserId   string `json:"user_id"`
	}{"", ""}
	err = srv.db.ReadSession(id, &dummyOneValid)
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
	_, errId := srv.db.CreateSession(user.UserId, user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", Api()+"/session/"+"id", nil)
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
	id, errId := srv.db.CreateSession(user.UserId, user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("GET", Api()+"/session/"+id, nil)
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
	id, errId := srv.db.CreateSession(user.UserId, user, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	req, err := http.NewRequest("DELETE", Api()+"/session/"+"id", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Errorf("Expected an error when deleteing non existing key but got StatusOK")
	}
	req, err = http.NewRequest("DELETE", Api()+"/session/"+id, nil)
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

func TestUpdateSession(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	userCreate := SessionAttributes{
		"userOne",
		"device_one",
		"",
		"",
		"",
		"userCreate@email.provider",
	}
	id, errId := srv.db.CreateSession(userCreate.UserId, userCreate, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	userUpdate := struct {
		SessionId string            `json:"session_id"`
		Session   SessionAttributes `json:"session_attributes"`
	}{
		id,
		SessionAttributes{
			"userOne",
			"device_one",
			"auth_code_one",
			"access_token_one",
			"refresh_token_one",
			"userCreate@email.provider",
		},
	}
	js, err := json.Marshal(userUpdate)

	req, err := http.NewRequest("PATCH", Api()+"/session", bytes.NewBuffer(js))
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Failed to patch a session. Got error code %d", w.Code)
	}
	var userPatched SessionAttributes
	err = json.NewDecoder(w.Body).Decode(&userPatched)
	if err != nil {
		t.Errorf("Failed to decode PATCH /session http response. Error: %s", err.Error())
	}
	if userUpdate.Session.UserId != userPatched.UserId ||
		userUpdate.Session.AuthenticationCode != userPatched.AuthenticationCode ||
		userUpdate.Session.AccessToken != userPatched.AccessToken ||
		userUpdate.Session.RefreshToken != userPatched.RefreshToken ||
		userUpdate.Session.UserEmail != userPatched.UserEmail ||
		userUpdate.Session.DeviceId != userPatched.DeviceId {

		t.Errorf("Response from the service does not match the requested update.")
	}

	mr.Close()
}
func TestUpdateSessionFailed(t *testing.T) {
	mr, srv, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	userCreate := SessionAttributes{
		"userOne",
		"device_one",
		"",
		"",
		"",
		"userCreate@email.provider",
	}
	id, errId := srv.db.CreateSession(userCreate.UserId, userCreate, 10)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	requestWithError := func(body interface{}, reason string) {
		if js, err := json.Marshal(body); err != nil {
			t.Error(err)
		} else {
			if req, err := http.NewRequest("PATCH", Api()+"/session", bytes.NewBuffer(js)); err != nil {
				t.Error(err)
			} else {
				w := httptest.NewRecorder()
				srv.ServeHTTP(w, req)
				if w.Code == http.StatusOK {
					t.Errorf("Expecting a request to fail due to " + reason + ", but request succeded.")
				}
			}
		}
	}
	userUpdate := struct {
		SessionId string            `json:"session_id"`
		Session   SessionAttributes `json:"session_attributes"`
	}{
		id,
		SessionAttributes{
			"userOne",
			"device_one",
			"auth_code_one",
			"access_token_one",
			"refresh_token_one",
			"userCreate@email.provider",
		},
	}
	if js, err := json.Marshal(userUpdate); err != nil {
		t.Error(err)
	} else {
		if req, err := http.NewRequest("PATCH", Api()+"/session", bytes.NewBuffer(js)); err != nil {
			t.Error(err)
		} else {
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("Expecting a request to succede, but request faileds.")
			}
		}
	}

	userUpdate.SessionId = "non-existing-session"
	requestWithError(userUpdate, "non-existing session id")

	userUpdate.SessionId = id
	userUpdate.Session.UserId = "unknown_user"
	requestWithError(userUpdate, "non-existing user id")

	var userUpdateNoBody = struct {
		SessionId string            `json:"session_id"`
		Session   map[string]string `json:"session_attributes"`
	}{
		SessionId: id,
		Session:   map[string]string{"": ""},
	}
	requestWithError(userUpdateNoBody, "empty request body")

	mr.Close()
}
