package routing

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

func Method(method string) ports.RequestMatcher {
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return strings.EqualFold(req.Method, method)
	})
}

func HeaderPresent(header string) ports.RequestMatcher {
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return req.Header.Get(header) != ""
	})
}

func Header(header, want string) ports.RequestMatcher {
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		values := req.Header.Values(header)
		if len(values) < 1 {
			return false
		}

		for _, value := range values {
			if strings.EqualFold(value, want) {
				return true
			}
		}

		return false
	})
}

func Path(path string) ports.RequestMatcher {
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return req.URL.Path == path
	})
}

func PathPattern(pattern string) (ports.RequestMatcher, error) {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return compiledPattern.MatchString(req.URL.Path)
	}), nil
}

func Query(key, value string) ports.RequestMatcher {
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return req.URL.Query().Get(key) == value
	})
}

func QueryPattern(key, pattern string) (ports.RequestMatcher, error) {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile pattern %s: %w", pattern, err)
	}
	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		return compiledPattern.MatchString(req.URL.Query().Get(key))
	}), nil
}

func JsonPath(path string, want any) (ports.RequestMatcher, error) {
	expression, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}

	return ports.RequestMatcherFunc(func(req *domain.IncomingRequest) bool {
		data, err := req.Body.Data()
		if err != nil {
			return false
		}

		parsed, err := oj.Parse(data)
		if err != nil {
			return false
		}

		for _, val := range expression.Get(parsed) {
			if reflect.DeepEqual(want, val) {
				return true
			}
		}

		return false
	}), nil
}
