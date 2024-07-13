package http

import (
	"log/slog"
	"net/http"

	"github.com/pb33f/libopenapi/renderer"
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
