package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/apple/pkl-go/pkl"

	"github.com/prskr/go-dito/core/ports"
	http2 "github.com/prskr/go-dito/handlers/http"
	"github.com/prskr/go-dito/infrastructure/config"
	"github.com/prskr/go-dito/infrastructure/grammar"
	"github.com/prskr/go-dito/infrastructure/routing"
)

var _ ports.SpecParser = (*PlainRuleSpecParser)(nil)

func NewPlainRuleSpecParser(spec *config.PlainRuleSpec) (*PlainRuleSpecParser, error) {
	var ruleSpecParser PlainRuleSpecParser

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for rule := range spec.Rules {
		slog.Info("Parsing DSL rule", slog.String("rule", rule))
		resp, err := grammar.Parse[grammar.ResponsePipeline](rule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule %s: %w", rule, err)
		}

		matcher, err := routing.ParseMatcher(resp.Filters())
		if err != nil {
			return nil, fmt.Errorf("failed to parse matcher %s: %w", rule, err)
		}

		responseProvider, err := routing.ParseResponseProvider(resp.Response, ports.CWD(os.DirFS(cwd)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse response provider %s: %w", rule, err)
		}

		ruleSpecParser.Handlers = append(ruleSpecParser.Handlers, http2.RulesRequestHandler{
			Matcher:          matcher,
			ResponseProvider: responseProvider,
		})
	}

	return &ruleSpecParser, nil
}

type PlainRuleSpecParser struct {
	Handlers []ports.RequestHandler
}

func (p PlainRuleSpecParser) Handler(ctx context.Context) (http.Handler, error) {
	return http2.RulesHandler{
		MaxContentLength: int64(config.Current().Server.RequestOptions.MaxBodySize.ToUnit(pkl.Bytes).Value),
		Handlers:         p.Handlers,
	}, nil
}
