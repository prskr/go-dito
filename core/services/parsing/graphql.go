package parsing

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/core/services/grammar"
	"github.com/prskr/go-dito/core/services/routing"
	httpHandlers "github.com/prskr/go-dito/handlers/http"
)

var _ ports.SpecParser = (*GraphQL)(nil)

type GraphQL struct {
	Schemes []string `json:"schemes"`
	Rules   []string `json:"rules"`
}

func (g GraphQL) Handler(ctx context.Context) (http.Handler, error) {
	sources := make([]*ast.Source, 0, len(g.Schemes))

	for _, schemaSrc := range g.Schemes {
		data, err := os.ReadFile(schemaSrc)
		if err != nil {
			return nil, err
		}

		sources = append(sources, &ast.Source{
			Name:    filepath.Base(schemaSrc),
			Input:   string(data),
			BuiltIn: false,
		})
	}

	handlers := make([]*httpHandlers.RulesRequestHandler, 0, len(g.Rules))

	schema, err := gqlparser.LoadSchema(sources...)
	if err != nil {
		return nil, err
	}

	parser := routing.GqlParser{Schema: schema}

	for _, rule := range g.Rules {
		slog.Info("Parsing GraphQL DSL rule", slog.String("rule", rule))
		resp, err := grammar.Parse[grammar.ResponsePipeline](rule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule %s: %w", rule, err)
		}

		matcher, err := parser.ParseMatchers(resp.Filters())
		if err != nil {
			return nil, fmt.Errorf("failed to parse matcher %s: %w", rule, err)
		}

		responseProvider, err := routing.ParseResponseProvider(resp.Response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response provider %s: %w", rule, err)
		}

		handlers = append(handlers, &httpHandlers.RulesRequestHandler{
			Matcher:          matcher,
			ResponseProvider: responseProvider,
		})
	}

	handler := httpHandlers.RulesHandler{
		Handlers: make([]ports.RequestHandler, len(handlers)),
	}

	for idx, h := range handlers {
		handler.Handlers[idx] = h
	}

	return handler, nil
}
