package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	if err := json.Unmarshal(inrec, &uInfoMap); err != nil {
		return uInfoMap
	}
	return uInfoMap
}
func objectFromMap(m map[string]interface{}, userInfo interface{}) error {
	inrec, _ := json.Marshal(m)
	json.Unmarshal(inrec, userInfo)

	return nil
}

func makeSessionKey(id string) string {
	return "session:" + id
}
func makeUserKey(id string) string {
	return "user:" + id
}

// Creates a session for the provided user id.
// If user id is not provided the function fails and no records in the database are created.
// Returns a session id string on success.
func (s *Service) createSession(userId string, sessionAttribs interface{}, expirationSec int) (string, error) {
	uInfoMap := objectToMap(sessionAttribs)

	rndStr, errRnd := s.randomString(47)
	if errRnd != nil {
		return "", errRnd
	}
	sessionId := makeSessionKey(rndStr)
	//TODO replace iteration through a map after a miniredis fix will be available
	for k, v := range uInfoMap {
		err := s.db.HSet(sessionId, k, v).Err()
		if err != nil {
			return "", err
		}
	}
	if expirationSec >= 0 {
		err := s.db.Expire(sessionId, time.Duration(expirationSec)*time.Second).Err()
		if err != nil {
			return "", nil
		}
	}
	if err := s.db.RPush(makeUserKey(userId), rndStr).Err(); err != nil {
		return "", err
	}
	return rndStr, nil
}

// Returns session attributes for the provided sessionId
// If sessionId does not exist returns error
// If none of the requested attributes found returns error
// If sessionId exists and at least one of the attributes exists returns error == nil
//	and fills the output object dest.
//	The attributes which do not exist are replaced by the empty string.
func (s *Service) readSession(sessionId string, dest interface{}) error {
	m := objectToMap(dest)
	keys := make([]string, len(m))
	values := make([]interface{}, len(m))
	var i = 0
	//TODO convert to HMGET
	for k, _ := range m {
		keys[i] = k
		value, err := s.db.HGet(makeSessionKey(sessionId), k).Result()
		if err != nil || value == "" {
			values[i] = nil
		} else {
			values[i] = value
		}
		i++
	}
	isNil := true
	for i := 0; i < len(keys); i++ {
		if values[i] != nil {
			isNil = false
			m[keys[i]] = values[i]
		} else {
			m[keys[i]] = ""
		}
	}
	if isNil {
		return errors.New("requested attributes not found")
	}
	return objectFromMap(m, dest)
}

// Deletes session key and related session attributes.
func (s *Service) deleteSession(sessionId string) error {
	session := makeSessionKey(sessionId)
	// Retrieve a user name associated with the session.
	// Remove a session id from the users's list of sessions.
	if user, err := s.db.HGet(session, "user_id").Result(); err != nil {
		return err
	} else {
		s.db.LRem(makeUserKey(user), 0, sessionId)
	}
	// Delete a session itself.
	if err := s.db.Del(session).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Service) updateSession(sessionId string, src interface{}, dest interface{}) error {
	sessionKey := makeSessionKey(sessionId)
	res, err := s.db.Exists(sessionKey).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New(fmt.Sprintf("Session \"%s\"Not found", sessionId))
	}
	m := objectToMap(src)
	pipe := s.db.TxPipeline()
	for k, v := range m {
		if err := pipe.HSet(sessionKey, k, v).Err(); err != nil {
			pipe.Discard()
			return err
		}
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := s.readSession(sessionId, dest); err != nil {
		return err
	}
	return nil
}
