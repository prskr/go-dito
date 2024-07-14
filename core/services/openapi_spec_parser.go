package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/renderer"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/prskr/go-dito/core/ports"
	http2 "github.com/prskr/go-dito/handlers/http"
	"github.com/prskr/go-dito/infrastructure/config"
	"github.com/prskr/go-dito/internal/maps"
)

var _ ports.SpecParser = (*OpenAPISpecParser)(nil)

const contentType = "application/json"

func NewOpenAPISpecParser(spec *config.OpenApiSpec) (*OpenAPISpecParser, error) {
	rawSchema, err := os.ReadFile(spec.SchemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", spec.SchemaPath, err)
	}

	specDocument, err := libopenapi.NewDocument(rawSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema file %s: %w", spec.SchemaPath, err)
	}

	return &OpenAPISpecParser{
		Spec:          specDocument,
		MockGenerator: renderer.NewMockGenerator(renderer.JSON),
	}, nil
}

type OpenAPISpecParser struct {
	Spec          libopenapi.Document
	MockGenerator *renderer.MockGenerator
}

func (o OpenAPISpecParser) Handler(ctx context.Context) (http.Handler, error) {
	mux := http.NewServeMux()

	v := o.Spec.GetVersion()

	switch {
	case strings.HasPrefix(v, "2"):
		if err := o.handleV2(ctx, mux); err != nil {
			return nil, err
		}

		return mux, nil

	case strings.HasPrefix(v, "3"):
		model, errs := o.Spec.BuildV3Model()
		if errs != nil {
			return nil, errors.Join(errs...)
		}

		if err := o.handleV3(ctx, mux, model); err != nil {
			return nil, err
		}

		schemaValidator := validator.NewValidatorFromV3Model(&model.Model)

		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			isValid, validationErrors := schemaValidator.ValidateHttpRequest(request)
			if !isValid {
				resp := struct {
					ValidationErrors []string
				}{}

				for _, err := range validationErrors {
					resp.ValidationErrors = append(resp.ValidationErrors, err.Error())
				}

				writer.Header().Set("Content-Type", contentType)
				writer.WriteHeader(http.StatusBadRequest)
				encoder := json.NewEncoder(writer)
				_ = encoder.Encode(resp)
			}

			mux.ServeHTTP(writer, request)
		}), nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", v)
	}
}

func (o OpenAPISpecParser) handleV3(ctx context.Context, mux *http.ServeMux, model *libopenapi.DocumentModel[v3.Document]) error {
	for path, ops := range maps.Iter(model.Model.Paths.PathItems) {
		logger := slog.Default().With(slog.String("api", model.Model.Info.Title), slog.String("path", path))

		for httpMethod, operation := range maps.Iter(ops.GetOperations()) {
			logger = logger.With(slog.String("http_method", httpMethod))

			pattern := fmt.Sprintf("%s %s", strings.ToUpper(httpMethod), path)

			response := operation.Responses.Default
			if response == nil {

				for rawStatus, responseValue := range maps.Iter(operation.Responses.Codes) {
					statusCode, err := strconv.ParseInt(rawStatus, 10, 32)
					if err != nil {
						logger.Warn("Failed to parse response status code", slog.String("status_code", rawStatus))
						continue
					}

					if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
						mediaType, present := responseValue.Content.Get(contentType)
						if !present {
							logger.Warn("No JSON response defined")
							continue
						}

						oasRulesCounter.Add(
							ctx,
							1,
							metric.WithAttributes(
								attribute.String("api", model.Model.Info.Title),
								attribute.String("version", model.Model.Version),
								attribute.StringSlice("tags", operation.Tags),
							),
						)

						if mediaType.Examples.IsZero() {
							logger.Info("Configuring mock handler")
							mockHandler := http2.OASSchemaMockHandler{
								MockGenerator: o.MockGenerator,
								Schema:        mediaType.Schema.Schema(),
								Status:        int(statusCode),
							}

							mux.Handle(pattern, otelhttp.WithRouteTag(pattern, mockHandler))
						} else {
							// TODO parse rule and configure handler
						}

						break
					}
				}
			}
		}
	}

	return nil
}

func (o OpenAPISpecParser) handleV2(ctx context.Context, mux *http.ServeMux) error {
	model, errs := o.Spec.BuildV2Model()
	if errs != nil {
		return errors.Join(errs...)
	}

	for path, ops := range maps.Iter(model.Model.Paths.PathItems) {
		logger := slog.Default().With(slog.String("api", model.Model.Info.Title), slog.String("path", path))

		for httpMethod, operation := range maps.Iter(ops.GetOperations()) {
			logger = logger.With(slog.String("http_method", httpMethod))

			pattern := fmt.Sprintf("%s %s", strings.ToUpper(httpMethod), path)

			response := operation.Responses.Default
			if response == nil {

				for rawStatus, responseValue := range maps.Iter(operation.Responses.Codes) {
					statusCode, err := strconv.ParseInt(rawStatus, 10, 32)
					if err != nil {
						logger.Warn("Failed to parse response status code", slog.String("status_code", rawStatus))
						continue
					}

					if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {

						oasRulesCounter.Add(
							ctx,
							1,
							metric.WithAttributes(
								attribute.String("api", model.Model.Info.Title),
								attribute.String("version", model.Model.Swagger),
								attribute.StringSlice("tags", operation.Tags),
							),
						)

						if responseValue.Examples == nil {
							logger.Info("Configuring mock handler")
							mockHandler := http2.OASSchemaMockHandler{
								MockGenerator: o.MockGenerator,
								Schema:        responseValue.Schema.Schema(),
								Status:        int(statusCode),
							}

							mux.Handle(pattern, otelhttp.WithRouteTag(pattern, mockHandler))
						} else {
							// TODO parse rule and configure handler
						}
					}
				}
			}
		}
	}

	return nil
}
