package routing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/prskr/go-dito/core/domain"
)

func TestGraphQlQuery_Matches(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	rawSchema, err := os.ReadFile(filepath.Join(cwd, "..", "..", "testdata", "star_wars_schema.graphql"))
	assert.NoError(t, err)

	schema := gqlparser.MustLoadSchema(&ast.Source{
		Name:    "star_wars_schema.graphql",
		Input:   string(rawSchema),
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
}`), 1<<31),
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
}`), 1<<31),
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
