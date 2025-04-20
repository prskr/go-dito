package grammar_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	grammar2 "github.com/prskr/go-dito/core/services/grammar"
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
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - string argument",
			rule:   `=> File("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Name:   "File",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - no argument",
			rule:   `=> NoContent()`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{Name: "NoContent"},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - no argument",
			rule:   `=> http.NoContent()`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "NoContent",
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - string argument",
			rule:   `=> http.File("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "File",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - int argument",
			rule:   `=> ReturnInt(1)`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Name:   "ReturnInt",
					Params: params(grammar2.Param{Int: grammar2.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response with module - int argument",
			rule:   `=> http.ReturnInt(1)`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(grammar2.Param{Int: grammar2.IntP(1)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response only - int argument, multiple digits",
			rule:   `=> ReturnInt(1337)`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Name:   "ReturnInt",
					Params: params(grammar2.Param{Int: grammar2.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - Response with Module - int argument, multiple digits",
			rule:   `=> http.ReturnInt(1337)`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "ReturnInt",
					Params: params(grammar2.Param{Int: grammar2.IntP(1337)}),
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - path pattern and terminator",
			rule:   `PathPattern(".*\\.(?i)png") => ReturnFile("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Name:   "ReturnFile",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
				FilterChain: &grammar2.Filters{
					Chain: []grammar2.Call{
						{
							Name:   "PathPattern",
							Params: params(grammar2.Param{String: grammar2.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - path pattern and terminator with Modules",
			rule:   `http.PathPattern(".*\\.(?i)png") => http.ReturnFile("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
				FilterChain: &grammar2.Filters{
					Chain: []grammar2.Call{
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(grammar2.Param{String: grammar2.StringP(`.*\.(?i)png`)}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - HTTP method, path pattern and terminator",
			rule:   `Method("GET") -> PathPattern("/index.html") => ReturnFile("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Name:   "ReturnFile",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
				FilterChain: &grammar2.Filters{
					Chain: []grammar2.Call{
						{
							Name:   "Method",
							Params: params(grammar2.Param{String: grammar2.StringP(http.MethodGet)}),
						},
						{
							Name:   "PathPattern",
							Params: params(grammar2.Param{String: grammar2.StringP("/index.html")}),
						},
					},
				},
			},
			wantErr: false,
		},
		parseTest[grammar2.ResponsePipeline]{
			name:   "ResponsePipeline - HTTP method, path pattern and terminator with modules",
			rule:   `http.Method("GET") -> http.PathPattern("/index.html") => http.ReturnFile("default.html")`,
			parser: grammar2.Parse[grammar2.ResponsePipeline],
			want: &grammar2.ResponsePipeline{
				Response: &grammar2.Call{
					Module: "http",
					Name:   "ReturnFile",
					Params: params(grammar2.Param{String: grammar2.StringP("default.html")}),
				},
				FilterChain: &grammar2.Filters{
					Chain: []grammar2.Call{
						{
							Module: "http",
							Name:   "Method",
							Params: params(grammar2.Param{String: grammar2.StringP(http.MethodGet)}),
						},
						{
							Module: "http",
							Name:   "PathPattern",
							Params: params(grammar2.Param{String: grammar2.StringP("/index.html")}),
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
