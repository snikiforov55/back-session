package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

const UserIdVarName = "user_id"
const SessionIdName = "session_id"
const SessionAttrName = "session_attributes"

func RandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
func ObjectToMap(userInfo interface{}) map[string]interface{} {
	var uInfoMap map[string]interface{}
	inrec, _ := json.Marshal(userInfo)
	if err := json.Unmarshal(inrec, &uInfoMap); err != nil {
		return uInfoMap
	}
	return uInfoMap
}
func ObjectFromMap(m map[string]interface{}, userInfo interface{}) error {
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

type Database interface {
	// Creates a session for the provided user id.
	// If user id is not provided the function fails and no records in the database are created.
	// Returns a session id string on success.
	createSession(userId string, sessionAttribs interface{}, expirationSec int) (string, error)

	// Returns session attributes for the provided sessionId
	// If sessionId does not exist returns error
	// If none of the requested attributes found returns error
	// If sessionId exists and at least one of the attributes exists returns error == nil
	//	and fills the output object dest.
	//	The attributes which do not exist are replaced by the empty string.
	readSession(sessionId string, dest interface{}) error

	// Deletes session key and related session attributes.
	deleteSession(sessionId string) error

	// Updates session attributes.
	// The userId cannot be changed. If it is provided in the payload
	// it should match the user_id provided when the session was created.
	// If it does not match the update fails.
	// The userId may be provided in the updatePayload but will be ignored in the updates.
	updateSession(sessionId string, updatePayload interface{}, dest interface{}) error
}
