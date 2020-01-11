package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func (s *Service) createSession(userInfo interface{}, expiration time.Duration) (string, error) {
	uInfo, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	rndStr, errRnd := s.randomString(47)
	if errRnd != nil {
		return "", errRnd
	}
	err = s.db.HMSet("session:"+rndStr, uInfo, expiration).Err()
	if err != nil {
		return "", err
	}
	return rndStr, nil
}
