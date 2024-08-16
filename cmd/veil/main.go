package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// GenerateInterfaceWithMethods generates a Go interface that includes all methods for a given struct type (from ast.TypeSpec).
func GenerateInterfaceWithMethods(file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Ensure the type is a struct.
	_, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return "", fmt.Errorf("%s is not a struct", typeSpec.Name.Name)
	}

	// Generate interface name based on the struct name.
	interfaceName := fmt.Sprintf("%sInterface", typeSpec.Name.Name)
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("type %s interface {\n", interfaceName))

	// Retrieve all methods associated with the struct.
	methods := GetMethodsForStruct(file, typeSpec.Name.Name)

	// Iterate over the methods and generate method signatures.
	for _, method := range methods {
		methodSignature := GenerateMethodSignature(method)
		builder.WriteString(fmt.Sprintf("\t%s\n", methodSignature))
	}

	builder.WriteString("}\n")

	return builder.String(), nil
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

// GenerateMethodSignature generates a method signature string for a given FuncDecl.
func GenerateMethodSignature(funcDecl *ast.FuncDecl) string {
	var builder strings.Builder
	builder.WriteString(funcDecl.Name.Name)

	// Handle the function parameters.
	builder.WriteString("(")
	for i, param := range funcDecl.Type.Params.List {
		for j, name := range param.Names {
			if i > 0 || j > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(name.Name + " " + getTypeAsString(param.Type))
		}
	}
	builder.WriteString(") ")

	// Handle the return values.
	if funcDecl.Type.Results != nil {
		builder.WriteString("(")
		for i, result := range funcDecl.Type.Results.List {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(getTypeAsString(result.Type))
		}
		builder.WriteString(")")
	}

	return builder.String()
}

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

// Testing
func main() {
	// Replace "your/project/path" with the actual path to your project
	// projectPath :=  "/Users/jstrohm/code/veil/cmd/veil"
	// fmt.Println(os.Environ())
	fileName := os.Getenv("GOFILE")
	pkgName := os.Getenv("GOPACKAGE")
	ifile := "impl_" + fileName

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing directory:", err)
		return
	}

	// Store the comments in the file.
	var lastComment string

	var builder strings.Builder
	builder.WriteString("package " + pkgName + "\n\n")

	ast.Inspect(file, func(n ast.Node) bool {
		// Check for comments first.
		if cg, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range cg.List {
				if strings.Contains(comment.Text, "d:service") {
					lastComment = comment.Text // Save the comment if it contains "d:service"
				}
			}
		}

		// We're looking for type specifications (struct declarations).
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if the type is a struct.
		if _, ok := typeSpec.Type.(*ast.StructType); ok {
			// If there's a "d:service" comment, associate it with the struct.
			if lastComment != "" {
				// Reset the last comment after it is used.
				lastComment = ""

				s, _ := GenerateInterfaceWithMethods(file, typeSpec)
				builder.WriteString("// Generated from " + fileName + "\n")
				builder.WriteString(s)
				builder.WriteString("\n")
			}
		}

		return true
	})

	os.WriteFile(ifile, []byte(builder.String()), 0644)

}
