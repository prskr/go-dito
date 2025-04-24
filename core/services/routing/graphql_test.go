package routing

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/prskr/go-dito/core/domain"
)

//go:embed testdata/star_wars_schema.graphql
var rawStarWarsSchema []byte

func TestGraphQlQuery_Matches(t *testing.T) {
	schema := gqlparser.MustLoadSchema(&ast.Source{
		Name:    "star_wars_schema.graphql",
		Input:   string(rawStarWarsSchema),
		BuiltIn: false,
	})

	tests := []struct {
		name  string
		query string
		req   *domain.IncomingRequest
		want  bool
	}{
		{
			name: "Simple query - expect equality",
			// language=graphql
			query: `query asdf {
	allFilms {
		films {
			title
			director
		}
	}
}`,
			req: domain.NewRequest(graphQLRequest(
				// language=graphql
				`query {
	allFilms {
		films {
			director
			title
		}
	}
}`)),
			want: true,
		},
		{
			name: "Nested query - expect equality",
			// language=graphql
			query: `query {
	allFilms {
		films {
			title
			director
			characterConnection {
				characters {
					homeworld {
						name
						id
					}
					birthYear
					id
				}
			}
		}
	}
}`,
			req: domain.NewRequest(graphQLRequest(
				// language=graphql
				`query {
	allFilms {
		films {
			director
			title
			characterConnection {
				characters {
					id
					birthYear
					homeworld {
						id
						name
					}
				}
			}
		}
	}
}`)),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GraphQlInlineQuery{
				RawQuery: tt.query,
			}

			assert.NoError(t, g.InjectSchema(schema))
			assert.Equalf(t, tt.want, g.Matches(tt.req), "Matches(%v)", tt.req)
		})
	}
}

func graphQLRequest(query string) *http.Request {
	replacer := strings.NewReplacer("\n", " ", "\t", "")
	return (&http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(fmt.Sprintf(`{"query": "%s"}`, replacer.Replace(query)))),
	}).WithContext(context.Background())
}
