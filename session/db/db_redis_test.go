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

func setupServer(funcRandString func(_ int) (string, error)) (*miniredis.Miniredis, db.Database, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}
	client := db.NewRedisClient(&db.RedisConfig{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}, funcRandString)
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
	mr, dbc, errSrv := setupServer(testRandomString)
	if errSrv != nil {
		t.Error(errSrv)
		return
	}
	defer mr.Close()
	user := session.SessionAttributes{
		UserId:   "userOne",
		DeviceId: "Netscape on Nokia phone",
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
}

func TestUpdateSessionDB(t *testing.T) {
	mr, dbc, errSrv := setupServer(testRandomString)
	if errSrv != nil {
		t.Error(errSrv)
		return
	}
	defer mr.Close()
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

}

func TestReadUserSessionsDB(t *testing.T) {
	mr, dbc, errSrv := setupServer(db.RandomString)
	if errSrv != nil {
		t.Error(errSrv)
		return
	}
	defer mr.Close()

	sessionInit := struct {
		UserId   string `json:"user_id"`
		DeviceId string `json:"device_id"`
	}{
		UserId: "userOne",
	}
	devices := []string{"computer", "laptop", "phone_1", "phone_2"}
	type (
		Session struct {
			SessionId string
			Device    string
		}
	)
	sessionIds := []Session{}
	for _, d := range devices {
		sessionInit.DeviceId = d
		id, err := dbc.CreateSession(sessionInit.UserId, sessionInit, 10)
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
		sessionIds = append(sessionIds, Session{id, d})
	}
	sessions, err := dbc.ReadUserSessions(sessionInit.UserId, []string{"device_id"})
	if err != nil {
		t.Error(err)
		return
	}
	if len(sessions) != len(devices) {
		t.Errorf("Number of returned sessions %d does not match the number of created sessions %d",
			len(sessions), len(devices))
		return
	}
	for _, session := range sessionIds {
		s, presents := sessions[session.SessionId]
		if !presents {
			t.Errorf("Created session ID %s was not returned by the ReadUserSessions", session)
			return
		}
		if dev := s["device_id"]; dev != session.Device {
			t.Errorf("Returned device id %s does not match the created device id %s", dev, session.Device)
			return
		}
	}
}
