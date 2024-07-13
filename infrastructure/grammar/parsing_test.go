package grammar_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prskr/go-dito/infrastructure/grammar"
)

type testCase interface {
	Run(t *testing.T)
	Name() string
}

type parseTest[T any] struct {
	name    string
	rule    string
	parser  func(rule string) (*T, error)
	want    any
	wantErr bool
}

func (pt parseTest[T]) Name() string {
	return pt.name
}

func (pt parseTest[T]) Run(t *testing.T) {
	t.Helper()
	t.Parallel()
	got, err := pt.parser(pt.rule)

	if (err != nil) != pt.wantErr {
		t.Errorf("pt.wantErr = %v but got error %v", pt.wantErr, err)
	}

	if pt.wantErr {
		return
	}

	assert.Equal(t, got, pt.want)
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - string argument",
			rule:   `=> File("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Name:   "File",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - no argument",
			rule:   `=> NoContent()`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{Name: "NoContent"},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - no argument",
			rule:   `=> http.NoContent()`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "NoContent",
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - string argument",
			rule:   `=> http.File("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "File",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - int argument",
			rule:   `=> ReturnInt(1)`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Name:   "ReturnInt",
					Params: params(grammar.Param{Int: grammar.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - int argument",
			rule:   `=> http.ReturnInt(1)`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(grammar.Param{Int: grammar.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - int argument, multiple digits",
			rule:   `=> ReturnInt(1337)`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Name:   "ReturnInt",
					Params: params(grammar.Param{Int: grammar.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - Response with Module - int argument, multiple digits",
			rule:   `=> http.ReturnInt(1337)`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(grammar.Param{Int: grammar.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - path pattern and terminator",
			rule:   `PathPattern(".*\\.(?i)png") => ReturnFile("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Name:   "ReturnFile",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
				FilterChain: &grammar.Filters{
					Chain: []grammar.Call{
						{
							Name:   "PathPattern",
							Params: params(grammar.Param{String: grammar.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - path pattern and terminator with Modules",
			rule:   `http.PathPattern(".*\\.(?i)png") => http.ReturnFile("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
				FilterChain: &grammar.Filters{
					Chain: []grammar.Call{
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(grammar.Param{String: grammar.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - HTTP method, path pattern and terminator",
			rule:   `Method("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Name:   "ReturnFile",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
				FilterChain: &grammar.Filters{
					Chain: []grammar.Call{
						{
							Name:   "Method",
							Params: params(grammar.Param{String: grammar.StringP(http.MethodGet)}),
						},
						{
							Name:   "PathPattern",
							Params: params(grammar.Param{String: grammar.StringP("/index.html")}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar.ResponsePipeline]{
			name:   "ResponsePipeline - HTTP method, path pattern and terminator with modules",
			rule:   `http.Method("GET") -> http.PathPattern("/index.html") => http.ReturnFile("default.html")`,
			parser: grammar.Parse[grammar.ResponsePipeline],
			want: &grammar.ResponsePipeline{
				Response: &grammar.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(grammar.Param{String: grammar.StringP("default.html")}),
				},
				FilterChain: &grammar.Filters{
					Chain: []grammar.Call{
						{
							Module: "http",
							Name:   "Method",
							Params: params(grammar.Param{String: grammar.StringP(http.MethodGet)}),
						},
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(grammar.Param{String: grammar.StringP("/index.html")}),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	//nolint:paralleltest // is actually called in Run function
	for _, tc := range tests {
		tt := tc
		t.Run(tt.Name(), tt.Run)
	}
}
