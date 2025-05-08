package routing

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/slices"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/infrastructure/logging"
)

var (
	ignoredSelections = []string{
		"__typename",
	}
	_              ports.RequestMatcher = (*GraphQlFileQuery)(nil)
	errSchemaIsNil                      = errors.New("schema is nil")
)

type graphQlMatcherBase struct {
	Schema *ast.Schema
	Query  *ast.QueryDocument
}

func (g graphQlMatcherBase) Matches(req *domain.IncomingRequest) bool {
	ctx, span := tracer.Start(req.Context(), "Matches")
	defer span.End()

	if g.Schema == nil {
		span.RecordError(errSchemaIsNil)
		return false
	}

	bodyReader, err := req.Body.Reader()
	if err != nil {
		span.RecordError(err)
		slog.WarnContext(ctx, "failed to read request body", logging.Error(err))
		return false
	}

	span.AddEvent("Read request body")

	decoder := json.NewDecoder(bodyReader)

	var graphqlBody struct {
		Query     string          `json:"query"`
		Variables json.RawMessage `json:"variables"`
	}

	if err := decoder.Decode(&graphqlBody); err != nil {
		span.RecordError(err)
		slog.WarnContext(ctx, "failed to parse request body", logging.Error(err))
	}

	span.AddEvent("Decoded Query from request body")
	span.SetAttributes(
		attribute.String("GraphQLQuery", graphqlBody.Query),
		attribute.String("GraphQLVariables", string(graphqlBody.Variables)))

	queryDoc, errList := gqlparser.LoadQuery(g.Schema, graphqlBody.Query)
	if errList != nil && len(errList.Unwrap()) > 0 {
		slog.WarnContext(ctx, "failed to load query", logging.Error(errList))
		return false
	}

	span.AddEvent("Parsed Query")

	return slices.EqualFunc(g.Query.Operations, queryDoc.Operations, compareOperation)
}

func GraphQlQueryOf(schema *ast.Schema, query string) ports.RequestMatcher {
	return &GraphQlInlineQuery{
		graphQlMatcherBase: graphQlMatcherBase{
			Schema: schema,
		},
		RawQuery: query,
	}
}

type GraphQlInlineQuery struct {
	graphQlMatcherBase
	RawQuery string
}

func (g *GraphQlInlineQuery) InjectSchema(schema *ast.Schema) error {
	g.Schema = schema

	var errList gqlerror.List
	g.Query, errList = gqlparser.LoadQuery(g.Schema, g.RawQuery)
	if len(errList) > 0 {
		return errList
	}

	return nil
}

var _ ports.RequestMatcher = (*GraphQlFileQuery)(nil)

func GraphQlQueryFrom(schema *ast.Schema, filePath string) ports.RequestMatcher {
	return &GraphQlFileQuery{
		graphQlMatcherBase: graphQlMatcherBase{
			Schema: schema,
		},
		FilePath: filePath,
	}
}

type GraphQlFileQuery struct {
	graphQlMatcherBase
	FilePath string
}

func compareOperation(op1, op2 *ast.OperationDefinition) bool {
	if op1.Operation != op2.Operation {
		return false
	}

	if !slices.EqualFunc(op1.VariableDefinitions, op2.VariableDefinitions, variableDefinitionEquals) {
		slog.Debug("GraphQL variable definitions did not match")
		return false
	}

	op1.SelectionSet = filterSelectionSet(op1.SelectionSet)
	op2.SelectionSet = filterSelectionSet(op2.SelectionSet)

	slices.SortFunc(op1.SelectionSet, selectionCompareTo)
	slices.SortFunc(op2.SelectionSet, selectionCompareTo)

	if !slices.EqualFunc(op1.SelectionSet, op2.SelectionSet, selectionEquals) {
		slog.Debug("GraphQL selection set did not match")
		return false
	}

	if !slices.EqualFunc(op1.Directives, op2.Directives, directivesEquals) {
		slog.Debug("GraphQL drives of operation did not match")
		return false
	}

	return true
}

func variableDefinitionEquals(def1 *ast.VariableDefinition, def2 *ast.VariableDefinition) bool {
	return strings.EqualFold(def1.Variable, def2.Variable) &&
		strings.EqualFold(def1.Type.Name(), def2.Type.Name())
}

func directivesEquals(dir1 *ast.Directive, dir2 *ast.Directive) bool {
	if dir1.Name != dir2.Name {
		return false
	}

	// Compare arguments
	if len(dir1.Arguments) != len(dir2.Arguments) {
		return false
	}
	for i, arg1 := range dir1.Arguments {
		arg2 := dir2.Arguments[i]
		if arg1.Name != arg2.Name || arg1.Value.Raw != arg2.Value.Raw {
			return false
		}
	}

	return true
}

func filterSelectionSet(selSet ast.SelectionSet) ast.SelectionSet {
	if len(selSet) == 0 {
		return nil
	}

	filteredSet := make(ast.SelectionSet, 0, len(selSet))
	for _, selItem := range selSet {
		switch sel := selItem.(type) {
		case *ast.Field:
			if _, found := slices.BinarySearch(ignoredSelections, sel.Name); found {
				continue
			}

			filteredSet = append(filteredSet, sel)
		}
	}

	return filteredSet
}

func selectionCompareTo(sel1 ast.Selection, sel2 ast.Selection) int {
	switch s1 := sel1.(type) {
	case *ast.Field:
		s2, ok := sel2.(*ast.Field)
		if !ok {
			return -1
		}

		if len(s1.SelectionSet) > 0 {
			s1.SelectionSet = filterSelectionSet(s1.SelectionSet)
			slices.SortFunc(s1.SelectionSet, selectionCompareTo)
		}

		if len(s2.SelectionSet) > 0 {
			s2.SelectionSet = filterSelectionSet(s2.SelectionSet)
			slices.SortFunc(s2.SelectionSet, selectionCompareTo)
		}

		return strings.Compare(s1.Name, s2.Name)
	}

	return 0
}

func selectionEquals(sel1 ast.Selection, sel2 ast.Selection) bool {
	switch s1 := sel1.(type) {
	case *ast.Field:
		s2, ok := sel2.(*ast.Field)
		if !ok {
			return false
		}

		return fieldsEquals(s1, s2)
	case *ast.FragmentSpread:
		s2, ok := sel2.(*ast.FragmentSpread)
		if !ok {
			return false
		}
		return s1.Name == s2.Name && slices.EqualFunc(s1.Directives, s2.Directives, directivesEquals)
	case *ast.InlineFragment:
		s2, ok := sel2.(*ast.InlineFragment)
		if !ok {
			return false
		}

		s1.SelectionSet = filterSelectionSet(s1.SelectionSet)
		s2.SelectionSet = filterSelectionSet(s2.SelectionSet)

		return s1.TypeCondition == s2.TypeCondition &&
			slices.EqualFunc(s1.Directives, s2.Directives, directivesEquals) &&
			slices.EqualFunc(s1.SelectionSet, s2.SelectionSet, selectionEquals)
	}
	return true
}

func fieldsEquals(field1 *ast.Field, field2 *ast.Field) bool {
	if !strings.EqualFold(field1.Name, field2.Name) {
		return false
	}

	field1.SelectionSet = filterSelectionSet(field1.SelectionSet)
	field2.SelectionSet = filterSelectionSet(field2.SelectionSet)

	if selections1 := len(field1.SelectionSet); selections1 > 0 {
		if selections2 := len(field2.SelectionSet); selections2 != selections1 {
			return false
		}

		slices.SortFunc(field1.SelectionSet, selectionCompareTo)
		slices.SortFunc(field2.SelectionSet, selectionCompareTo)

		if !slices.EqualFunc(field1.SelectionSet, field2.SelectionSet, selectionEquals) {
			return false
		}
	}

	return true
}
