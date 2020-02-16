package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

const UserIdAttr = "user_id"
const CreateDateAttr = "create_date"
const UpdateDateAttr = "update_date"

//const SessionIdName = "session_id"
//const SessionAttrName = "session_attributes"

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
	if err := json.Unmarshal(inrec, userInfo); err != nil {
		return err
	}
	return nil
}
func MakeSessionKey(id string) string {
	return "session:" + id
}
func MakeUserKey(id string) string {
	return "user:" + id
}

type Database interface {
	// Returns a random string. Can be customized externally.
	RandomString() (string, error)
	// Creates a session for the provided user id.
	// If user id is not provided the function fails and no records in the database are created.
	// Returns a session id string on success.
	// By default the timestamp with a creation date is added to the session.
	// The session attributes contain default attributes:
	// { user_id, create_time, update_time }
	//
	CreateSession(userId string, sessionAttribs interface{}, expirationSec int) (string, error)

	// Returns session attributes for the provided sessionId
	// If sessionId does not exist returns error
	// If none of the requested attributes found returns error
	// If sessionId exists and at least one of the attributes exists returns error == nil
	//	and fills the output object dest.
	//	The attributes which do not exist are replaced by the empty string.
	ReadSession(sessionId string, dest interface{}) error

	// Read all sessions belonging to the user.
	// Fills destination map with requested session attributes, where a map
	// key is a session id, and value is a map of attribute_key:attribute_value.
	// Attribute Key and Attribute Value are strings.
	// For example:
	// {"session 1" : {"device_id": "a device one"},
	//  "session 2" : {"device_id": "another device"},
	// }
	// If non-existing attribute is requested, the attribute is ignored.
	ReadUserSessions(userId string, attributes []string) (map[string]map[string]string, error)

	// Deletes session key and related session attributes.
	DeleteSession(sessionId string) error

	// Updates session attributes.
	// The userId cannot be changed. If it is provided in the payload
	// it should match the user_id provided when the session was created.
	// If it does not match the update fails.
	// The userId may be provided in the updatePayload but will be ignored in the updates.
	UpdateSession(sessionId string, updatePayload interface{}, dest interface{}) error
}
