package db

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	db           redis.Cmdable
	randomString func(int) (string, error)
}
type RedisConfig struct {
	RedisHost     string `json:"redis_host"`
	RedisPort     string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDb       int    `json:"redis_db"`
}

func NewRedisClient(config *RedisConfig, randomString func(int) (string, error)) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDb,       // use default DB
	})
	return &Redis{client, randomString}
}

func (s *Redis) RandomString() (string, error) {
	return s.randomString(47)
}

// Creates a session for the provided user id.
// If user id is not provided the function fails and no records in the database are created.
// Returns a session id string on success.
func (s *Redis) CreateSession(userId string, sessionAttribs interface{}, expirationSec int) (string, error) {
	uInfoMap := ObjectToMap(sessionAttribs)

	rndStr, errRnd := s.RandomString()
	if errRnd != nil {
		return "", errRnd
	}
	sessionId := MakeSessionKey(rndStr)
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
	if err := s.db.RPush(MakeUserKey(userId), rndStr).Err(); err != nil {
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
func (s *Redis) ReadSession(sessionId string, dest interface{}) error {
	m := ObjectToMap(dest)
	keys := make([]string, len(m))
	values := make([]interface{}, len(m))
	var i = 0
	//TODO convert to HMGET
	for k, _ := range m {
		keys[i] = k
		value, err := s.db.HGet(MakeSessionKey(sessionId), k).Result()
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
	return ObjectFromMap(m, dest)
}

func (s *Redis) ReadUserSessions(userId string, attributes []string) (map[string]map[string]string, error) {
	userSessions, err := s.db.LRange(MakeUserKey(userId), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	getAttributes := func(sessionId string, attribs []string) (map[string]string, error) {
		output := map[string]string{}
		for _, reqAttrib := range attribs {
			attr, err := s.db.HGet(sessionId, reqAttrib).Result()
			if err != nil {
				return nil, err
			}
			output[reqAttrib] = attr
		}
		return output, nil
	}
	userSessionsContent := map[string]map[string]string{}
	for _, session := range userSessions {
		attributes, err := getAttributes(MakeSessionKey(session), attributes)
		if err != nil {
			return nil, err
		}
		userSessionsContent[session] = attributes
	}
	return userSessionsContent, nil
}

// Deletes session key and related session attributes.
func (s *Redis) DeleteSession(sessionId string) error {
	session := MakeSessionKey(sessionId)
	// Retrieve a user name associated with the session.
	// Remove a session id from the users's list of sessions.
	if user, err := s.db.HGet(session, UserIdVarName).Result(); err != nil {
		return err
	} else {
		s.db.LRem(MakeUserKey(user), 0, sessionId)
	}
	// Delete a session itself.
	if err := s.db.Del(session).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Redis) UpdateSession(sessionId string, src interface{}, dest interface{}) error {
	sessionKey := MakeSessionKey(sessionId)
	res, err := s.db.Exists(sessionKey).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New(fmt.Sprintf("SessionId \"%s\"Not found", sessionId))
	}
	m := ObjectToMap(src)
	// If user_id is provided in the payload check if it matches the provided user id.
	// If user id does not match the database return error.
	if uid, ok := m[UserIdVarName]; ok {
		uidDb, err := s.db.HGet(sessionKey, UserIdVarName).Result()
		if err != nil {
			return errors.New(fmt.Sprintf("SessionId \"%s\" does not have a user id.", sessionId))
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
		if k == UserIdVarName {
			continue
		}
		if err := pipe.HSet(sessionKey, k, v).Err(); err != nil {
			if dis_err := pipe.Discard(); dis_err != nil {
				return errors.New(err.Error() + " followed by " + dis_err.Error())
			}
			return err
		}
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := s.ReadSession(sessionId, dest); err != nil {
		return err
	}
	return nil
}
