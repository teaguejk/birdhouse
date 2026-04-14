package constants

type ContextKey string

const (
	DeviceIDKey     ContextKey = "deviceID"
	DeviceActiveKey ContextKey = "deviceActive"
	OAuthClaimsKey  ContextKey = "oauthClaims"
)
