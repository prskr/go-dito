package ports

import "github.com/vektah/gqlparser/v2/ast"

type CwdInjectable interface {
	InjectCwd(cwd CWD)
}

type SchemaInjectable interface {
	InjectSchema(schema *ast.Schema) error
}
