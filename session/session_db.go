package session

import (
	"errors"
	"fmt"
	"github.com/snikiforov55/back-session/session/db"
	"time"
)

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
	uInfoMap := db.ObjectToMap(sessionAttribs)

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

// Returns session attributes for the provided SessionId
// If SessionId does not exist returns error
// If none of the requested attributes found returns error
// If SessionId exists and at least one of the attributes exists returns error == nil
//	and fills the output object dest.
//	The attributes which do not exist are replaced by the empty string.
func (s *Service) readSession(sessionId string, dest interface{}) error {
	m := db.ObjectToMap(dest)
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
	return db.ObjectFromMap(m, dest)
}

// Deletes session key and related session attributes.
func (s *Service) deleteSession(sessionId string) error {
	session := makeSessionKey(sessionId)
	// Retrieve a user name associated with the session.
	// Remove a session id from the users's list of sessions.
	if user, err := s.db.HGet(session, db.UserIdVarName).Result(); err != nil {
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
	m := db.ObjectToMap(src)
	// If user_id is provided in the payload check if it matches the provided user id.
	// If user id does not match the database return error.
	if uid, ok := m[db.UserIdVarName]; ok {
		uidDb, err := s.db.HGet(sessionKey, db.UserIdVarName).Result()
		if err != nil {
			return errors.New(fmt.Sprintf("Session \"%s\" does not have a user id.", sessionId))
		}
		if uidDb != uid {
			return errors.New(fmt.Sprintf("User id \"%s\" stored in the session \"%s\" does not "+
				"match a requested user id \"%s\".", uidDb, sessionId, uid))
		}
	}
	// Continue with the update.
	pipe := s.db.TxPipeline()
	for k, v := range m {
		// Do not set user id
		if k == db.UserIdVarName {
			continue
		}
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
