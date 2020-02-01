package session

const DefaultSessionExpirationSec = 3600 * 24 * 365

type SessionAttributes struct {
	UserId             string `json:"user_id"`
	DeviceId           string `json:"device_id"`
	AuthenticationCode string `json:"auth_code"`
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	UserEmail          string `json:"user_email"`
}
