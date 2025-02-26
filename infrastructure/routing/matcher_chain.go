package routing

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

var (
	_ ports.RequestMatcher   = (RequestMatcherChain)(nil)
	_ ports.SchemaInjectable = (RequestMatcherChain)(nil)
	_ ports.CwdInjectable    = (RequestMatcherChain)(nil)
)

type RequestMatcherChain []ports.RequestMatcher

func (r RequestMatcherChain) InjectCwd(cwd ports.CWD) {
	for _, m := range r {
		if injectable, ok := m.(ports.CwdInjectable); ok {
			injectable.InjectCwd(cwd)
		}
	}
}

func (r RequestMatcherChain) InjectSchema(schema *ast.Schema) error {
	for _, m := range r {
		if injectable, ok := m.(ports.SchemaInjectable); ok {
			if err := injectable.InjectSchema(schema); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r RequestMatcherChain) Matches(req *domain.IncomingRequest) bool {
	for _, m := range r {
		if !m.Matches(req) {
			return false
		}
	}

	return true
}
