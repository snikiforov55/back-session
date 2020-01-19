## Session Manager Service
The purpose of the Session Manager service is to keep sessions of the user.
The user can have multiple sessions simultaneously open from different devices.
The session is a key-value pair. 
The Key is a generated unique session ID (GUID string).
The session attributes contain:
```
{
    user_id : String   // Required. Unique user identifier in the system.
    auth_code : String // Authentication code received from the authentication/authorization server
    access_token : String // Access token issued by the Auth server
    refresh_token : String // Refresh token issued by the Auth server
    exiration_time : Int // Expiration time in seconds for the access token.   
    user_email: String // User e-mail
    device_id: String // Required. Device identification string from which user has logged in
}
```
The session short attributes contain:
```
{
    user_id : String
    sessions: [
        {
            device_id: String
            session_id: String
        }
    ]
}
```

User can view owned sessions.
User can close one session.
User can close all sessions.

Session management service is implemented in GO.
The session database is redis.

### Session API
|    Command and Path    |  Parameters          | Return               | Description                                                
|------------------------|----------------------|----------------------|--------------------------
| GET    /session        | user_id              | session_short(id, device) | Get All Sessions for current user
| GET    /session/{id}   |                      | session_attributes(user_id, access_key, access_token, refresh_token) | Get Session details
| POST   /session        | user_id              | session_id                  | Create new session
| PATCH  /session/{id}   | session_attributes   | session_attributes  | Update auth info for the session. If new tokens are provided.
| DELETE /session/{id}   |                      | session_id           | Delete session
| DELETE /session        | user_id              |             | Delete all sessions for a user
