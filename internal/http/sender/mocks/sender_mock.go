package sender_mocks

import (
	httpsender "festwrap/internal/http/sender"

	"github.com/stretchr/testify/mock"
)

type HTTPSenderMock struct {
	mock.Mock
}

func (s *HTTPSenderMock) Send(options httpsender.HTTPRequestOptions) (*[]byte, error) {
	args := s.Called(options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	bytes := args.Get(0).(*[]byte)
	return bytes, args.Error(1)
}
