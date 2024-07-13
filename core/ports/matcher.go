package ports

import (
	"net/http"

	"github.com/prskr/go-dito/core/domain"
)

type RequestMatcher interface {
	Matches(req *domain.IncomingRequest) bool
}

type RequestMatcherFunc func(req *domain.IncomingRequest) bool

func (f RequestMatcherFunc) Matches(req *domain.IncomingRequest) bool {
	return f(req)
}

type ResponseProvider interface {
	Apply(writer http.ResponseWriter)
}

type ResponseProviderFunc func(writer http.ResponseWriter)

func (f ResponseProviderFunc) Apply(writer http.ResponseWriter) {
	f(writer)
}

type RequestHandler interface {
	Handle(writer http.ResponseWriter, ir *domain.IncomingRequest) (handled bool)
}
