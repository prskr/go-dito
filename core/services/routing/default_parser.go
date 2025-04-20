package routing

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/core/services/grammar"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func init() {
	slices.Sort(httpMethods)
}

type DefaultParser struct{}

func (p DefaultParser) ParseMatchers(filters []grammar.Call) (ports.RequestMatcher, error) {
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

func (DefaultParser) ParseMatcher(filterCall grammar.Call) (ports.RequestMatcher, error) {
	switch filterCall.Signature() {
	case "http.method(string)":
		method, _ := filterCall.Params[0].AsString()
		_, found := slices.BinarySearch(httpMethods, method)
		if !found {
			return nil, fmt.Errorf("method '%s' is not a valid HTTP method", method)
		}

		return Method(method), nil
	case "http.headerpresent(string)":
		headerName, _ := filterCall.Params[0].AsString()

		return HeaderPresent(headerName), nil
	case "http.header(string,string)":
		headerName, _ := filterCall.Params[0].AsString()
		headerValue, _ := filterCall.Params[1].AsString()

		return Header(headerName, headerValue), nil
	case "http.path(string)":
		path, _ := filterCall.Params[0].AsString()

		return Path(path), nil
	case "http.pathpattern(string)":
		pathPattern, _ := filterCall.Params[0].AsString()

		return PathPattern(pathPattern)
	case "http.query(string,string)":
		queryKey, _ := filterCall.Params[0].AsString()
		queryValue, _ := filterCall.Params[1].AsString()

		return Query(queryKey, queryValue), nil
	case "http.querypattern(string,string)":
		queryKey, _ := filterCall.Params[0].AsString()
		queryValuePattern, _ := filterCall.Params[1].AsString()

		return QueryPattern(queryKey, queryValuePattern)
	case "http.jsonpath(string,string)":
		jsonPath, _ := filterCall.Params[0].AsString()
		jsonPathMatcher, err := JsonPath(jsonPath, filterCall.Params[1].Value())
		if err != nil {
			return nil, err
		}

		return jsonPathMatcher, nil
	default:
		return nil, fmt.Errorf("unknown filter call: %s", filterCall.String())
	}
}
