package httpsender

type HTTPRequestSender interface {
	Send(options HTTPRequestOptions) (*[]byte, error)
}

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

type HTTPRequestOptions struct {
	body               []byte
	url                string
	method             Method
	headers            map[string]string
	expectedStatusCode int
}

func (o *HTTPRequestOptions) GetBody() []byte {
	return o.body
}

func (o *HTTPRequestOptions) GetUrl() string {
	return o.url
}

func (o *HTTPRequestOptions) GetMethod() Method {
	return o.method
}

func (o *HTTPRequestOptions) GetHeaders() map[string]string {
	return o.headers
}

func (o *HTTPRequestOptions) GetExpectedStatusCode() int {
	return o.expectedStatusCode
}

func (o *HTTPRequestOptions) SetUrl(url string) {
	o.url = url
}

func (o *HTTPRequestOptions) SetBody(body []byte) {
	o.body = body
}

func (o *HTTPRequestOptions) SetHeaders(headers map[string]string) {
	o.headers = headers
}

func NewHTTPRequestOptions(url string, method Method, expectedStatusCode int) HTTPRequestOptions {
	return HTTPRequestOptions{url: url, body: nil, method: method, expectedStatusCode: expectedStatusCode}
}
