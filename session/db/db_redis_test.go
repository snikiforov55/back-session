package db_test

import (
	"github.com/alicebob/miniredis"
	"github.com/snikiforov55/back-session/session"
	"github.com/snikiforov55/back-session/session/db"
	"testing"
)

func testRandomString(_ int) (string, error) {
	return "SESSION_01", nil
}

func setupServer() (*miniredis.Miniredis, db.Database, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}
	client := db.NewRedisClient(&db.RedisConfig{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}, testRandomString)
	return mr, client, nil
}

func TestRandomString(t *testing.T) {
	str, err := db.RandomString(37)
	if err != nil {
		t.Error(err)
	}
	if len(str) < 37 {
		t.Errorf("Unexpected random string length. Waiting for >= 47 got %d", len(str))
	}
}

func TestUpdateAuthInfo(t *testing.T) {

}

func TestReadSessionDbFailures(t *testing.T) {
	mr, dbc, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	user := session.SessionAttributes{
		"userOne",
		"Netscape on Nokia phone",
		"",
		"",
		"",
		"",
	}
	id, errId := dbc.CreateSession(user.UserId, user, 3600)
	if errId != nil {
		t.Errorf("%s", errId.Error())
	}
	err := dbc.ReadSession("id", &user)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid session id but it didn't.")
	}
	dummy := struct {
		Whatever string `json:"whatever"`
	}{""}
	err = dbc.ReadSession(id, &dummy)
	if err == nil {
		t.Errorf("Expecting readSession to fail due to invalid structure but it didn't.")
	}
	dummyOneValid := struct {
		Whatever string `json:"whatever"`
		UserId   string `json:"user_id"`
	}{"", ""}
	err = dbc.ReadSession(id, &dummyOneValid)
	if err != nil {
		t.Errorf("Expecting readSession not to fail but it did.")
	}
	mr.Close()
}

func TestUpdateSessionDB(t *testing.T) {
	mr, dbc, errSrv := setupServer()
	if errSrv != nil {
		t.Error(errSrv)
	}
	sessionInit := struct {
		UserId string `json:"user_id"`
	}{
		"userOne",
	}
	id, errId := dbc.CreateSession(sessionInit.UserId, sessionInit, 10)
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
	err := dbc.UpdateSession(id, sessionUpd, &sessionRes)
	if err != nil {
		t.Errorf("Update session failed %s", err.Error())
	}
	if sessionRes.UserId != sessionUpd.UserId || sessionRes.AccessToken != sessionUpd.AccessToken {
		t.Errorf("Returned session does not match the update %s != %s", sessionUpd, sessionRes)
	}
	if err := dbc.UpdateSession("SessionId", sessionUpd, &sessionRes); err == nil {
		t.Error("Expecting a call fail due to invalid session ID but it didn't")
	}
	mr.Close()
}
