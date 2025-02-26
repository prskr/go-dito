package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prskr/go-dito/core/ports"
	handlers "github.com/prskr/go-dito/handlers/http"
	"github.com/prskr/go-dito/infrastructure/config"
	"github.com/prskr/go-dito/infrastructure/grammar"
	"github.com/prskr/go-dito/infrastructure/routing"
)

var (
	_ ports.SpecParser    = (*PlainRuleSpecParser)(nil)
	_ ports.CwdInjectable = (*PlainRuleSpecParser)(nil)
)

func NewPlainRuleSpecParser(spec *config.PlainRuleSpec) (*PlainRuleSpecParser, error) {
	var ruleSpecParser PlainRuleSpecParser

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

		responseProvider, err := routing.ParseResponseProvider(resp.Response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response provider %s: %w", rule, err)
		}

		ruleSpecParser.Handlers = append(ruleSpecParser.Handlers, handlers.RulesRequestHandler{
			Matcher:          matcher,
			ResponseProvider: responseProvider,
		})
	}

	return &ruleSpecParser, nil
}

type PlainRuleSpecParser struct {
	Handlers []ports.RequestHandler
}

func (p PlainRuleSpecParser) InjectCwd(cwd ports.CWD) {
	for _, handler := range p.Handlers {
		handler.(handlers.RulesRequestHandler).Init(cwd)
	}
}

func (p PlainRuleSpecParser) Handler(context.Context) (http.Handler, error) {
	maxBodySize := config.Current().Server.RequestOptions.MaxBodySize
	bodySizeLimit := int64(maxBodySize.Unit) * int64(maxBodySize.Value)

	return handlers.RulesHandler{
		MaxContentLength: bodySizeLimit,
		Handlers:         p.Handlers,
	}, nil
}
