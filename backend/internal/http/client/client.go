package httpclient

import (
	"net/http"
)

type HTTPClient interface {
	Send(request *http.Request) (*http.Response, error)
}
