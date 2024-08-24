package main

import (
	"fmt"
	"go/ast"
	"strings"
)

// getTypeAsString converts an ast.Expr (field type or function parameter/return type) to its string representation.
func getTypeAsString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		// For types like `time.Time`, where a package selector is used.
		return fmt.Sprintf("%s.%s", getTypeAsString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		// Handle pointer types.
		return fmt.Sprintf("*%s", getTypeAsString(t.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", getTypeAsString(t.Elt))
	case *ast.FuncType:
		// Handle function types (rare case for field types).
		var params []string
		for _, param := range t.Params.List {
			params = append(params, getTypeAsString(param.Type))
		}
		var results []string
		if t.Results != nil {
			for _, result := range t.Results.List {
				results = append(results, getTypeAsString(result.Type))
			}
		}
		return fmt.Sprintf("func(%s) (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
	default:
		return "unknown"
	}
}

type Builder struct {
	strings.Builder
}

func (b *Builder) Sprintf(format string, a ...any) {
	b.WriteString(fmt.Sprintf(format, a...))
}

func UppercaseFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}
