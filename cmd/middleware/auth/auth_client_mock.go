package auth

import "github.com/stretchr/testify/mock"

type AuthClientMock struct {
	mock.Mock
}

func (s *AuthClientMock) GetAccessToken() (string, error) {
	args := s.Called()
	if args.Get(0) == "" {
		return "", args.Error(1)
	}

	bytes := args.Get(0).(string)
	return bytes, args.Error(1)
}
