package routing

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/core/services/grammar"
)

type GqlParser struct {
	DefaultParser
	Schema *ast.Schema
}

func (p GqlParser) ParseMatchers(filters []grammar.Call) (ports.RequestMatcher, error) {
	compiledFilters := make([]ports.RequestMatcher, 0, len(filters))
	for _, filterCall := range filters {
		matcher, err := p.ParseMatcher(filterCall)
		if err != nil {
			return nil, err
		}
		compiledFilters = append(compiledFilters, matcher)
	}

	return RequestMatcherChain(compiledFilters), nil
}

func (p GqlParser) ParseMatcher(filterCall grammar.Call) (ports.RequestMatcher, error) {
	switch filterCall.Signature() {
	case "graphql.query(string)":
		gqlQuery, _ := filterCall.Params[0].AsString()
		return GraphQlQueryOf(p.Schema, gqlQuery), nil
	case "graphql.queryfromfile(string)":
		filePath, _ := filterCall.Params[0].AsString()
		return GraphQlQueryFrom(p.Schema, filePath), nil
	default:
		return p.DefaultParser.ParseMatcher(filterCall)
	}
}
