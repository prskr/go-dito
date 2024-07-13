package http

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

var _ ports.RequestHandler = (*RulesRequestHandler)(nil)

type RulesRequestHandler struct {
	Matcher          ports.RequestMatcher
	ResponseProvider ports.ResponseProvider
}

func (r RulesRequestHandler) Handle(writer http.ResponseWriter, ir *domain.IncomingRequest) (handled bool) {
	_, span := tracer.Start(ir.Original.Context(), "MatchRequestWithRule")
	defer span.End()

	defer func() {
		span.SetAttributes(attribute.Bool("matched", handled))
	}()

	if r.Matcher.Matches(ir) {
		r.ResponseProvider.Apply(writer)
		return true
	}

	return false
}
