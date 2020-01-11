package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"
)

func randomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
func objectToMap(userInfo interface{}) map[string]interface{} {
	var uInfoMap map[string]interface{}
	inrec, _ := json.Marshal(userInfo)
	json.Unmarshal(inrec, &uInfoMap)
	return uInfoMap
}
func objectFromMap(m map[string]interface{}, userInfo interface{}) error {
	inrec, _ := json.Marshal(m)
	json.Unmarshal(inrec, userInfo)

	return nil
}

func (s *Service) createSession(userInfo interface{}, expiration time.Duration) (string, error) {
	uInfoMap := objectToMap(userInfo)

	rndStr, errRnd := s.randomString(47)
	if errRnd != nil {
		return "", errRnd
	}
	sessionId := "session:" + rndStr
	//TODO replace iteration through a map after a miniredis fix will be available
	for k, v := range uInfoMap {
		err := s.db.HSet(sessionId, k, v).Err()
		if err != nil {
			return "", err
		}
	}
	err := s.db.Expire(sessionId, expiration).Err()
	if err != nil {
		return "", nil
	}
	return rndStr, nil
}

func (s *Service) readSession(sessionId string, dest *UserInfo) error {
	m := objectToMap(*dest)

	keys := make([]string, reflect.ValueOf(*dest).NumField())
	values := make([]string, reflect.ValueOf(*dest).NumField())
	var i = 0
	for k, _ := range m {
		keys[i] = k
		value, err := s.db.HGet("session:"+sessionId, k).Result()
		if err != nil {
			return err
		}
		values[i] = value
		i++
	}
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = values[i]
	}
	return objectFromMap(m, dest)
}
