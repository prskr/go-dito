package grammar

import (
	"errors"
	"fmt"
	"strings"
)

var ErrTypeMismatch = errors.New("param has a different type")

// ResponsePipeline describes how requests are handled that expect single response
// e.g. HTTP or DNS requests
// A ResponsePipeline is defined as an optional chain of Filters like: filter1() -> filter2()
// and a Response which determines how the request should be handled e.g. http.Status(204)
// a full chain might look like so: GET() -> Header("Accept", "application/json") -> http.Status(200).
type ResponsePipeline struct {
	FilterChain *Filters `parser:"@@*"`
	Response    *Call    `parser:"'=' '>' @@"`
}

func (p *ResponsePipeline) Filters() []Call {
	if p.FilterChain != nil {
		return p.FilterChain.Chain
	}
	return nil
}

type Filters struct {
	Chain []Call `parser:"@@ ('-' '>' @@)*"`
}

type Call struct {
	Module string  `parser:"(@Ident'.')?"`
	Name   string  `parser:"@Ident"`
	Params []Param `parser:"'(' @@? ( ',' @@ )*')'"`
}

func (c Call) Signature() string {
	types := make([]string, 0, len(c.Params))
	for _, param := range c.Params {
		types = append(types, param.Type())
	}

	if c.Module == "" {
		return strings.ToLower(fmt.Sprintf("%s(%s)", c.Name, strings.Join(types, ",")))
	}

	return strings.ToLower(fmt.Sprintf("%s.%s(%s)", c.Module, c.Name, strings.Join(types, ",")))
}

func (c Call) String() string {
	params := make([]string, 0, len(c.Params))
	for _, param := range c.Params {
		params = append(params, fmt.Sprintf("%v", param.Value()))
	}
	return fmt.Sprintf("%s.%s(%s)", c.Module, c.Name, strings.Join(params, ","))
}

type Param struct {
	String *string  `parser:"@String | @RawString"`
	Int    *int     `parser:"| @Int"`
	Float  *float64 `parser:"| @Float"`
}

func (p Param) AsString() (string, error) {
	if p.String == nil {
		return "", fmt.Errorf("string is nil %w", ErrTypeMismatch)
	}

	return *p.String, nil
}

func (p Param) AsInt() (int, error) {
	if p.Int == nil {
		return 0, fmt.Errorf("int is nil %w", ErrTypeMismatch)
	}

	return *p.Int, nil
}

func (p Param) AsFloat() (float64, error) {
	if p.Float == nil {
		return 0, fmt.Errorf("float is nil %w", ErrTypeMismatch)
	}

	return *p.Float, nil
}

func (p Param) Value() any {
	if p.String != nil {
		return *p.String
	}

	if p.Int != nil {
		return *p.Int
	}

	if p.Float != nil {
		return *p.Float
	}

	return nil
}

func (p Param) Type() string {
	if p.String != nil {
		return "string"
	}

	if p.Int != nil {
		return "int"
	}

	if p.Float != nil {
		return "float"
	}

	return "nil"
}
