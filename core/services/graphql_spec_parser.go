package services

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
	handlers "github.com/prskr/go-dito/handlers/http"
	"github.com/prskr/go-dito/infrastructure/config"
	"github.com/prskr/go-dito/infrastructure/grammar"
	"github.com/prskr/go-dito/infrastructure/routing"
)

var (
	_ ports.SpecParser    = (*GraphQLSchemaSpec)(nil)
	_ ports.CwdInjectable = (*GraphQLSchemaSpec)(nil)
)

func NewGraphQLSchemaSpec(spec *config.GraphQlSpec) (schemaSpec *GraphQLSchemaSpec, err error) {
	sources := make([]*ast.Source, 0, len(spec.Schemas))

	for schemaSrc := range spec.Schemas {
		data, err := os.ReadFile(schemaSrc.Path)
		if err != nil {
			return nil, err
		}

		sources = append(sources, &ast.Source{
			Name:    filepath.Base(schemaSrc.Path),
			Input:   string(data),
			BuiltIn: schemaSrc.BuiltIn,
		})
	}

	schemaSpec = new(GraphQLSchemaSpec)

	schemaSpec.Schema, err = gqlparser.LoadSchema(sources...)
	if err != nil {
		return nil, err
	}

	for rule := range spec.Rules {
		slog.Info("Parsing GraphQL DSL rule", slog.String("rule", rule))
		resp, err := grammar.Parse[grammar.ResponsePipeline](rule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule %s: %w", rule, err)
		}

		matcher, err := routing.ParseMatcher(resp.Filters())
		if err != nil {
			return nil, fmt.Errorf("failed to parse matcher %s: %w", rule, err)
		}

		responseProvider, err := routing.ParseResponseProvider(resp.Response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response provider %s: %w", rule, err)
		}

		schemaSpec.Handlers = append(schemaSpec.Handlers, &handlers.RulesRequestHandler{
			Matcher:          matcher,
			ResponseProvider: responseProvider,
		})
	}

	return schemaSpec, nil
}

type GraphQLSchemaSpec struct {
	Schema   *ast.Schema
	Handlers []ports.RequestHandler
}

func (g *GraphQLSchemaSpec) Handler(context.Context) (http.Handler, error) {
	for _, handler := range g.Handlers {
		if injectable, ok := handler.(*handlers.RulesRequestHandler).Matcher.(ports.SchemaInjectable); ok {
			if err := injectable.InjectSchema(g.Schema); err != nil {
				return nil, err
			}
		}
	}

	maxBodySize := config.Current().Server.RequestOptions.MaxBodySize
	bodySizeLimit := int64(maxBodySize.Unit) * int64(maxBodySize.Value)

	return handlers.RulesHandler{
		MaxContentLength: bodySizeLimit,
		Handlers:         g.Handlers,
	}, nil
}

func (g *GraphQLSchemaSpec) InjectCwd(cwd ports.CWD) {
	for _, handler := range g.Handlers {
		handler.(*handlers.RulesRequestHandler).Init(cwd)
	}
}
