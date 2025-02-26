package services

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/prskr/go-dito/core/ports"
)

var (
	ErrNotMatchingParser = errors.New("spec type does not match parser")
	ErrUnknownParserType = errors.New("unknown parser type")
	DefaultRegistry      = make(ParserRegistry)
)

func init() {
	RegisterAt(&DefaultRegistry, NewOpenAPISpecParser)
	RegisterAt(&DefaultRegistry, NewPlainRuleSpecParser)
	RegisterAt(&DefaultRegistry, NewGraphQLSchemaSpec)
}

type ParserRegistry map[reflect.Type]func(spec any) (ports.SpecParser, error)

func (r *ParserRegistry) ParserFor(spec any) (ports.SpecParser, error) {
	provider, ok := (*r)[reflect.TypeOf(spec)]
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrUnknownParserType, spec)
	}

	return provider(spec)
}

func RegisterAt[TSpec any, TParser ports.SpecParser](registry *ParserRegistry, provider func(spec TSpec) (TParser, error)) {
	var s TSpec

	(*registry)[reflect.TypeOf(s)] = func(spec any) (ports.SpecParser, error) {
		concreteSpec, ok := spec.(TSpec)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrNotMatchingParser, spec)
		}

		return provider(concreteSpec)
	}
}
