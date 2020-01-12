## Session Manager Service
The purpose of the Session Manager service is to keep sessions of the user.
The user can have multiple sessions simultaneously open from different devices.
The session is a key-value pair. 
The Key is a generated unique session ID (GUID string).
The session value contains:
* user id
* authentication code
* access token
* refresh token
* user e-mail
* device identification string

User can view owned sessions.
User can close one session.
User can close all sessions.

Session management service is implemented in GO.
The session database is redis.

### Session API
|    Command and Path    |  Parameters | Return               | Description                                                
|------------------------|-------------|----------------------|--------------------------
| GET    /session        | user_id     | session_details(id, device) | Get All Sessions for current user
| GET    /session/{id}   |             | session_detail(user_id, access_key, access_token, refresh_token) | Get Session details
| POST   /session        | user_id     | session_id  | Create new session
| PUT    /session/{id}   | auth_info   | session_id  | Update auth info for the session. If new tokens are provided.
| DELETE /session/{id}   |             |             | Delete session
| DELETE /session        |             |             | Delete all sessions
