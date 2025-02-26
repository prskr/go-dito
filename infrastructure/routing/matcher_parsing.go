package routing

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/infrastructure/grammar"
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

func ParseMatcher(filters []grammar.Call) (ports.RequestMatcher, error) {
	compiledFilters := make([]ports.RequestMatcher, 0, len(filters))
	for _, filterCall := range filters {
		switch filterCall.Signature() {
		case "http.method(string)":
			method, _ := filterCall.Params[0].AsString()
			_, found := slices.BinarySearch(httpMethods, method)
			if !found {
				return nil, fmt.Errorf("method '%s' is not a valid HTTP method", method)
			}

			compiledFilters = append(compiledFilters, Method(method))
		case "http.headerpresent(string)":
			headerName, _ := filterCall.Params[0].AsString()

			compiledFilters = append(compiledFilters, HeaderPresent(headerName))
		case "http.header(string,string)":
			headerName, _ := filterCall.Params[0].AsString()
			headerValue, _ := filterCall.Params[1].AsString()

			compiledFilters = append(compiledFilters, Header(headerName, headerValue))
		case "http.path(string)":
			path, _ := filterCall.Params[0].AsString()

			compiledFilters = append(compiledFilters, Path(path))
		case "http.pathpattern(string)":
			pathPattern, _ := filterCall.Params[0].AsString()

			pathPatternMatcher, err := PathPattern(pathPattern)
			if err != nil {
				return nil, err
			}

			compiledFilters = append(compiledFilters, pathPatternMatcher)
		case "http.query(string,string)":
			queryKey, _ := filterCall.Params[0].AsString()
			queryValue, _ := filterCall.Params[1].AsString()

			compiledFilters = append(compiledFilters, Query(queryKey, queryValue))
		case "http.querypattern(string,string)":
			queryKey, _ := filterCall.Params[0].AsString()
			queryValuePattern, _ := filterCall.Params[1].AsString()

			queryPatternMatcher, err := QueryPattern(queryKey, queryValuePattern)
			if err != nil {
				return nil, err
			}

			compiledFilters = append(compiledFilters, queryPatternMatcher)
		case "http.jsonpath(string, string)":
			jsonPath, _ := filterCall.Params[0].AsString()
			jsonPathMatcher, err := JsonPath(jsonPath, filterCall.Params[1].Value())
			if err != nil {
				return nil, err
			}

			compiledFilters = append(compiledFilters, jsonPathMatcher)
		case "graphql.query(string)":
			gqlQuery, _ := filterCall.Params[0].AsString()

			compiledFilters = append(compiledFilters, GraphQlQueryOf(gqlQuery))
		case "graphql.queryfromfile(string)":
			filePath, _ := filterCall.Params[0].AsString()

			compiledFilters = append(compiledFilters, GraphQlQueryFrom(filePath))
		default:
			return nil, fmt.Errorf("unknown filter call: %s", filterCall.String())
		}
	}

	return RequestMatcherChain(compiledFilters), nil
}
