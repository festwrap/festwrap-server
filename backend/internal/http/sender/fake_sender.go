package httpsender

type FakeHTTPSender struct {
	sendArgs HTTPRequestOptions
	response *[]byte
	err      error
}

func (s *FakeHTTPSender) GetSendArgs() HTTPRequestOptions {
	return s.sendArgs
}

func (s *FakeHTTPSender) SetResponse(response *[]byte) {
	s.response = response
}

func (s *FakeHTTPSender) SetError(err error) {
	s.err = err
}

func (s *FakeHTTPSender) Send(options HTTPRequestOptions) (*[]byte, error) {
	s.sendArgs = options

	if s.err != nil {
		return nil, s.err
	} else {
		return s.response, nil
	}
}

func NewFakeHTTPSender() FakeHTTPSender {
	return FakeHTTPSender{}
}
