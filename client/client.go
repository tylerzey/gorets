package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/jpfielding/gofilters/filter"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// const DefaultTimeout int = 300000

// standard http header names
const (
	UserAgent   string = "User-Agent"
	Accept      string = "Accept"
	ContentType string = "Content-Type"
	WWWAuth     string = "Www-Authenticate"
	WWWAuthResp string = "Authorization"
)

// rets http header names
const (
	RETSVersion   string = "RETS-Version"
	RETSSessionID string = "RETS-Session-ID"
	RETSRequestID string = "RETS-Request-ID"
	RETSUAAuth    string = "RETS-UA-Authorization"
)

// TODO create a Session interface with a Requester and a reset to clear state and pass that in

// Requester implmenters should not assume any order of ops
type Requester func(ctx context.Context, req *http.Request) (*http.Response, error)

// DefaultSession configures the default rets session
func DefaultSession(user, pwd, userAgent, userAgentPw, retsVersion string, transport http.RoundTripper) (Requester, error) {
	if transport == nil {
		transport = http.DefaultTransport
	}

	client := http.Client{
		Transport: transport,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	// 4) send the request
	request := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return ctxhttp.Do(ctx, &client, req)
	}
	// 3) www auth
	wwwAuth := (&WWWAuthTransport{
		Requester: request,
		Username:  user,
		Password:  pwd,
	}).Request
	// 2) apply ua auth headers per request, if there is a pwd
	uaAuth := (&UserAgentAuthentication{
		Requester:         wwwAuth,
		UserAgent:         userAgent,
		UserAgentPassword: userAgentPw,
		GetRETSVersion:    CreateRETSVersioner(retsVersion),
		GetSessionID:      CreateSessionIDer(client.Jar),
	}).Request
	// 1) apply default headers first (outermost wrapping)
	headers := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		req.Header.Set(UserAgent, userAgent)
		req.Header.Set(RETSVersion, retsVersion)
		req.Header.Set(Accept, "*/*")
		return uaAuth(ctx, req)
	}
	return headers, nil
}

// DefaultXMLDecoder the variable used to set a selected charset
var DefaultXMLDecoder = CreateXMLDecoder

// CreateXMLDecoder decodes xml using the given the header if needed
func CreateXMLDecoder(input io.Reader, strict bool) *xml.Decoder {
	// drop any chars that will blow up the xml decoder and replace with a space
	input = filter.NewReader(input, filter.XML10Filter(filter.SpaceChar))
	decoder := xml.NewDecoder(input)
	decoder.Strict = strict
	// this only gets used when a proper xml header is used
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
}

// DefaultReEncodeReader allows overriding the re-encoding operation
var DefaultReEncodeReader = ReEncodeReader

// ReEncodeReader re-encodes a reader based on the http content type provided
func ReEncodeReader(input io.ReadCloser, contentType string) io.ReadCloser {
	if e, _, _ := charset.DetermineEncoding([]byte{}, contentType); e != encoding.Nop {
		type closer struct {
			io.Reader
			io.Closer
		}
		tr := transform.NewReader(input, e.NewDecoder())
		return closer{tr, input}
	}
	return input
}
