package routing

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/exp/slices"

	"github.com/prskr/go-dito/core/domain"
	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/infrastructure/logging"
)

var (
	_ ports.RequestMatcher   = (*GraphQlFileQuery)(nil)
	_ ports.SchemaInjectable = (*GraphQlInlineQuery)(nil)
)

type graphQlMatcherBase struct {
	Schema *ast.Schema
	Query  *ast.QueryDocument
}

func (g graphQlMatcherBase) Matches(req *domain.IncomingRequest) bool {
	if g.Schema == nil {
		return false
	}

	bodyReader, err := req.Body.Reader()
	if err != nil {
		slog.WarnContext(req.Context(), "failed to read request body", logging.Error(err))
		return false
	}

	decoder := json.NewDecoder(bodyReader)

	var graphqlBody struct {
		Query string `json:"query"`
	}

	if err := decoder.Decode(&graphqlBody); err != nil {
		slog.WarnContext(req.Context(), "failed to parse request body", logging.Error(err))
	}

	queryDoc, errList := gqlparser.LoadQuery(g.Schema, graphqlBody.Query)
	if errList != nil && len(errList.Unwrap()) > 0 {
		slog.WarnContext(req.Context(), "failed to load query", logging.Error(errList))
		return false
	}

	return slices.EqualFunc(g.Query.Operations, queryDoc.Operations, compareOperation)
}

func GraphQlQueryOf(query string) ports.RequestMatcher {
	return &GraphQlInlineQuery{RawQuery: query}
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

var (
	_ ports.RequestMatcher   = (*GraphQlFileQuery)(nil)
	_ ports.SchemaInjectable = (*GraphQlFileQuery)(nil)
	_ ports.CwdInjectable    = (*GraphQlFileQuery)(nil)
)

func GraphQlQueryFrom(filePath string) ports.RequestMatcher {
	return &GraphQlFileQuery{
		FilePath: filePath,
	}
}

type GraphQlFileQuery struct {
	graphQlMatcherBase
	CWD      ports.CWD
	FilePath string
}

func (g *GraphQlFileQuery) InjectCwd(cwd ports.CWD) {
	g.CWD = cwd
}

func (g *GraphQlFileQuery) InjectSchema(schema *ast.Schema) error {
	g.Schema = schema

	rawQuery, err := fs.ReadFile(g.CWD, g.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", g.FilePath, err)
	}

	var errList gqlerror.List
	g.Query, errList = gqlparser.LoadQuery(g.Schema, string(rawQuery))
	if len(errList) > 0 {
		return errList
	}

	return nil
}

func compareOperation(op1, op2 *ast.OperationDefinition) bool {
	if op1.Operation != op2.Operation {
		return false
	}

	if !slices.EqualFunc(op1.VariableDefinitions, op2.VariableDefinitions, variableDefinitionEquals) {
		return false
	}

	slices.SortFunc(op1.SelectionSet, selectionCompareTo)
	slices.SortFunc(op2.SelectionSet, selectionCompareTo)

	if !slices.EqualFunc(op1.SelectionSet, op2.SelectionSet, selectionEquals) {
		return false
	}

	if !slices.EqualFunc(op1.Directives, op2.Directives, directivesEquals) {
		return false
	}

	return true
}

func variableDefinitionEquals(def1 *ast.VariableDefinition, def2 *ast.VariableDefinition) bool {
	return def1.Definition == def2.Definition
}

func directivesEquals(dir1 *ast.Directive, dir2 *ast.Directive) bool {
	if dir1.Name != dir2.Name {
		return false
	}

	return true
}

func selectionCompareTo(sel1 ast.Selection, sel2 ast.Selection) int {
	switch s1 := sel1.(type) {
	case *ast.Field:
		s2, ok := sel2.(*ast.Field)
		if !ok {
			return -1
		}

		if len(s1.SelectionSet) > 0 {
			slices.SortFunc(s1.SelectionSet, selectionCompareTo)
		}

		if len(s2.SelectionSet) > 0 {
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
	}
	return true
}

func fieldsEquals(field1 *ast.Field, field2 *ast.Field) bool {
	if !strings.EqualFold(field1.Name, field2.Name) {
		return false
	}

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
