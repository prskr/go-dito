package http

import (
	"net/http"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

var _ http.Handler = (*RulesHandler)(nil)

type RulesHandler struct {
	MaxContentLength int64
	Handlers         []ports.RequestHandler
}

func (r RulesHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx, span := tracer.Start(request.Context(), "HandleRequestWithRulesDSL")
	defer span.End()

	request = request.WithContext(ctx)

	ir := domain.NewRequest(request, r.MaxContentLength)

	for _, h := range r.Handlers {
		if handled := h.Handle(writer, ir); handled {
			return
		}
	}

	span.AddEvent("NoRuleMatched")

	http.NotFound(writer, request)
}
