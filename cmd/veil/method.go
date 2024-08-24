package main

import (
	"go/ast"
	"strings"
)

// GenerateMethodSignature generates a method signature string for a given FuncDecl.
func GenerateMethodSignature(funcDecl *ast.FuncDecl) string {
	var builder strings.Builder
	builder.WriteString(funcDecl.Name.Name)

	if len(funcDecl.Type.Params.List) == 0 {
		return ""
	}

	// Handle the function parameters.
	builder.WriteString("(")
	for i, param := range funcDecl.Type.Params.List {
		for j, name := range param.Names {
			if i > 0 || j > 0 {
				builder.WriteString(", ")
			}
			tas := getTypeAsString(param.Type)
			if i == 0 && tas != "context.Context" {
				return ""
			}
			builder.WriteString(name.Name + " " + tas)
		}
	}
	builder.WriteString(") ")

	// Handle the return values.
	if funcDecl.Type.Results != nil {
		builder.WriteString("(")
		errorTypeFound := false
		for i, result := range funcDecl.Type.Results.List {
			if i > 0 {
				builder.WriteString(", ")
			}
			tas := getTypeAsString(result.Type)
			if tas == "error" {
				errorTypeFound = true
			}
			builder.WriteString(tas)
		}
		builder.WriteString(")")
		if !errorTypeFound {
			return ""
		}
	}

	return builder.String()
}

// GetMethodsForStruct retrieves methods for a given struct type from the parsed AST.
func GetMethodsForStruct(file *ast.File, structName string) []*ast.FuncDecl {
	var methods []*ast.FuncDecl

	// Loop through declarations in the file.
	for _, decl := range file.Decls {
		// We are only interested in function declarations.
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if the function has a receiver.
		if funcDecl.Recv != nil {
			for _, receiver := range funcDecl.Recv.List {
				// We expect the receiver to be of type `*ast.StarExpr` for pointer receivers or `*ast.Ident` for value receivers.
				switch expr := receiver.Type.(type) {
				case *ast.Ident:
					// Check if the receiver type matches the struct name.
					if expr.Name == structName {
						methods = append(methods, funcDecl)
					}
				case *ast.StarExpr:
					// If the receiver is a pointer, the struct type is inside the StarExpr.
					if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == structName {
						methods = append(methods, funcDecl)
					}
				}
			}
		}
	}

	return methods
}
