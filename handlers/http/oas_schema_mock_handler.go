package http

import (
	"log/slog"
	"math/rand/v2"
	"net/http"

	"github.com/pb33f/libopenapi/renderer"
	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
)

var _ http.Handler = (*OASSchemaMockHandler)(nil)

type OASSchemaMockHandler struct {
	MockGenerator *renderer.MockGenerator
	Schema        any
	Status        int
}

func (h OASSchemaMockHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	_, generateMockSpan := tracer.Start(req.Context(), "GenerateMock")
	raw, err := h.MockGenerator.GenerateMock(h.Schema, "")
	if err != nil {
		generateMockSpan.RecordError(err)
		http.Error(writer, "Failed to generate mock", http.StatusInternalServerError)
		return
	}
	generateMockSpan.End()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(h.Status)

	if _, err := writer.Write(raw); err != nil {
		slog.Warn("Failed to write mock response")
	}
}

var _ http.Handler = (*OASSchemaExampleHandler)(nil)

type OASSchemaExampleHandler struct {
	Handlers       []ports.RequestHandler
	FallbackStatus int
	FallbackValues [][]byte
}

func (o OASSchemaExampleHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx, returnExampleSpan := tracer.Start(req.Context(), "ReturnExample")
	defer returnExampleSpan.End()

	req = req.WithContext(ctx)
	ir := domain.NewRequest(req)

	for _, handler := range o.Handlers {
		if handled := handler.Handle(writer, ir); handled {
			return
		}
	}

	returnExampleSpan.AddEvent("NoRuleMatched")

	if fallbackValueCount := len(o.FallbackValues); fallbackValueCount != 0 {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(o.FallbackStatus)

		_, _ = writer.Write(o.FallbackValues[rand.N(len(o.FallbackValues))])
		return
	}

	returnExampleSpan.AddEvent("NoFallbackValue")

	http.NotFound(writer, req)
}
