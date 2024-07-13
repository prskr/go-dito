package routing

import (
	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

var _ ports.RequestMatcher = (RequestMatcherChain)(nil)

type RequestMatcherChain []ports.RequestMatcher

func (r RequestMatcherChain) Matches(req *domain.IncomingRequest) bool {
	for _, m := range r {
		if !m.Matches(req) {
			return false
		}
	}

	return true
}
