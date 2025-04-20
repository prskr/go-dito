package parsing

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/core/services/grammar"
	"github.com/prskr/go-dito/core/services/routing"
	httpHandlers "github.com/prskr/go-dito/handlers/http"
)

var _ ports.SpecParser = (*Plain)(nil)

type Plain struct {
	Rules []string `json:"rules"`
}

func (p Plain) Handler(context.Context) (http.Handler, error) {

	var handlers []ports.RequestHandler

	var parser routing.DefaultParser

	for _, rule := range p.Rules {
		slog.Info("Parsing DSL rule", slog.String("rule", rule))
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

		handlers = append(handlers, httpHandlers.RulesRequestHandler{
			Matcher:          matcher,
			ResponseProvider: responseProvider,
		})
	}

	return nil, nil
}
