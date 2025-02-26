package domain

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
)

func NewRequest(orig *http.Request, maxContentLength int64) *IncomingRequest {
	return &IncomingRequest{
		Original:         orig,
		MaxContentLength: maxContentLength,
		Method:           orig.Method,
		URL:              orig.URL,
		Header:           orig.Header,
		Body: RequestBody{
			Original: orig,
		},
		Host:       orig.Host,
		RemoteAddr: orig.RemoteAddr,
	}
}

type RequestBody struct {
	readOnce         sync.Once
	data             []byte
	readErr          error
	MaxContentLength int64
	Original         *http.Request
}

func (b *RequestBody) Reader() (io.Reader, error) {
	b.readBody()

	if b.readErr != nil {
		return nil, b.readErr
	}

	return bytes.NewReader(b.data), nil
}

func (b *RequestBody) Data() ([]byte, error) {
	b.readBody()

	return b.data, b.readErr
}

func (b *RequestBody) readBody() {
	b.readOnce.Do(func() {
		b.data, b.readErr = io.ReadAll(b.Original.Body)
	})
}

type IncomingRequest struct {
	Original *http.Request

	MaxContentLength int64
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	// For client requests, an empty string means GET.
	Method string

	// URL specifies either the URI being requested (for server
	// requests) or the URL to access (for client requests).
	//
	// For server requests, the URL is parsed from the URI
	// supplied on the Request-Line as stored in RequestURI.  For
	// most requests, fields other than Path and RawQuery will be
	// empty. (See RFC 7230, Section 5.3)
	//
	// For client requests, the URL's Host specifies the server to
	// connect to, while the Request's Host field optionally
	// specifies the Host header value to send in the HTTP
	// request.
	URL *url.URL

	// Header contains the request header fields either received
	// by the server or to be sent by the client.
	//
	// If a server received a request with header lines,
	//
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"Bar", "two"},
	//	}
	//
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	//
	// HTTP defines that header names are case-insensitive. The
	// request parser implements this by using CanonicalHeaderKey,
	// making the first character and any characters following a
	// hyphen uppercase and the rest lowercase.
	//
	// For client requests, certain headers such as Content-Length
	// and Connection are automatically written when needed and
	// values in Header may be ignored. See the documentation
	// for the Request.Write method.
	Header http.Header

	// Body is the request's body.
	//
	// For client requests, a nil body means the request has no
	// body, such as a GET request. The HTTP Client's Transport
	// is responsible for calling the Close method.
	//
	// For server requests, the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	//
	// Body must allow Read to be called concurrently with Close.
	// In particular, calling Close should unblock a Read waiting
	// for input.
	Body RequestBody

	// For server requests, Host specifies the host on which the
	// URL is sought. For HTTP/1 (per RFC 7230, section 5.4), this
	// is either the value of the "Host" header or the host name
	// given in the URL itself. For HTTP/2, it is the value of the
	// ":authority" pseudo-header field.
	// It may be of the form "host:port". For international domain
	// names, Host may be in Punycode or Unicode form. Use
	// golang.org/x/net/idna to convert it to either format if
	// needed.
	// To prevent DNS rebinding attacks, server Handlers should
	// validate that the Host header has a value for which the
	// Handler considers itself authoritative. The included
	// ServeMux supports patterns registered to particular host
	// names and thus protects its registered Handlers.
	//
	// For client requests, Host optionally overrides the Host
	// header to send. If empty, the Request.Write method uses
	// the value of URL.Host. Host may contain an international
	// domain name.
	Host string

	// Form contains the parsed form data, including both the URL
	// field's query parameters and the PATCH, POST, or PUT form data.
	// This field is only available after ParseForm is called.
	// The HTTP client ignores Form and uses Body instead.
	Form url.Values

	// PostForm contains the parsed form data from PATCH, POST
	// or PUT body parameters.
	//
	// This field is only available after ParseForm is called.
	// The HTTP client ignores PostForm and uses Body instead.
	PostForm url.Values

	// MultipartForm is the parsed multipart form, including file uploads.
	// This field is only available after ParseMultipartForm is called.
	// The HTTP client ignores MultipartForm and uses Body instead.
	MultipartForm *multipart.Form

	// RemoteAddr allows HTTP servers and other software to record
	// the network address that sent the request, usually for
	// logging. This field is not filled in by ReadRequest and
	// has no defined format. The HTTP server in this package
	// sets RemoteAddr to an "IP:port" address before invoking a
	// handler.
	// This field is ignored by the HTTP client.
	RemoteAddr string
}

func (i *IncomingRequest) ParseMultipartForm() error {
	if err := i.Original.ParseMultipartForm(i.MaxContentLength); err != nil {
		return err
	}

	i.Form = i.Original.Form
	i.PostForm = i.Original.PostForm
	i.MultipartForm = i.Original.MultipartForm

	return nil
}

func (i *IncomingRequest) ParseForm() error {
	if err := i.Original.ParseForm(); err != nil {
		return err
	}

	i.Form = i.Original.Form
	i.PostForm = i.Original.PostForm

	return nil
}

func (i *IncomingRequest) Context() context.Context {
	return i.Original.Context()
}
