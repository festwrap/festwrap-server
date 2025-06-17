package auth

type AuthClient interface {
	GetAccessToken() (string, error)
}
