package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type collect struct {
	fileName   string
	projectDir string
	config     Config
}

func GetDataForGoFile(fileName string, config Config) (VeilData, error) {
	c := &collect{
		fileName:   fileName,
		config:     config,
		projectDir: filepath.Dir(fileName),
	}
	return c.get()
}

func containsMethod(arr []Method, value Method) bool {
	for _, v := range arr {
		if v.Name == value.Name {
			return true
		}
	}
	return false
}

func (c *collect) get() (VeilData, error) {
	fmt.Println("Getting Data for: ", c.fileName)

	data := VeilData{
		Filename:    c.fileName,
		PackageName: "main",
		Structs:     []Struct{},
		Types:       []string{},
	}

	// Let's do basic parse of the file, and find tags we care about
	taggedStructs, err := c.parseForTags()
	if err != nil {
		return data, err
	}

	// For all taggedStructs, look for embedded structs
	fmt.Println("Project Directory: ", c.projectDir)

	allStructs := map[string]Struct{}
	types := map[string]any{}

	err = filepath.Walk(c.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fmt.Println("Parsing: ", path)

		for _, holder := range taggedStructs {
			ts := holder.typeSpec
			structName := ts.Name.Name
			s, ok := allStructs[structName]
			if !ok {
				s = Struct{
					Name:    structName,
					Methods: []Method{},
				}
			}

			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("error parsing directory: %w", err)
			}

			methods, tempTypes := c.getAllMethods(astFile, structName, types)
			types = tempTypes
			s.Methods = append(s.Methods, methods...)

			for _, st := range holder.embedded {
				methods, tempTypes := c.getAllMethods(astFile, st, types)
				types = tempTypes

				for _, method := range methods {
					if !containsMethod(s.Methods, method) {
						s.Methods = append(s.Methods, method)
					}
				}
			}

			allStructs[structName] = s

		}

		return nil
	})
	if err != nil {
		return data, errors.Wrap(err, "Error walking paths")
	}

	// Now build the rest of the structure
	for _, v := range allStructs {
		data.Structs = append(data.Structs, v)
	}
	data.Types = extractTypes(types)

	return data, nil
}

func (c *collect) getAllMethods(astFile *ast.File, name string, types map[string]any) ([]Method, map[string]any) {
	// Generate method data structure
	methods := []Method{}

	for _, method := range GetMethodsForStruct(astFile, name) {
		mdata := Method{}
		mdata.Name = method.Name.Name

		args := []Arg{}

		for i, param := range method.Type.Params.List {
			for _, name := range param.Names {
				tas := getTypeAsString(param.Type)
				types[tas] = true
				if i == 0 {
					if tas != "context.Context" {
						goto skip
					}
					continue
				}

				args = append(args, Arg{
					Name: name.Name,
					Type: tas,
				})

			}
		}
	skip:
		mdata.Args = args

		if method.Type.Results != nil {
			errorTypeFound := false
			for _, result := range method.Type.Results.List {
				tas := getTypeAsString(result.Type)
				types[tas] = true
				if tas == "error" {
					errorTypeFound = true
				}
				mdata.Returns = append(mdata.Returns, tas)
			}
			if !errorTypeFound {
				mdata.Returns = []string{}
				goto skip2
			}
		}
	skip2:
		methods = append(methods, mdata)
	}

	return methods, types
}

type TypeHolder struct {
	typeSpec *ast.TypeSpec
	embedded []string
}

func (c *collect) parseForTags() ([]TypeHolder, error) {
	structs := []TypeHolder{}

	astFile, err := parser.ParseFile(token.NewFileSet(), c.fileName, nil, parser.ParseComments)
	if err != nil {
		return structs, fmt.Errorf("error parsing file: %w", err)
	}

	// Store the comments in the file.
	var lastComment string

	ast.Inspect(astFile, func(n ast.Node) bool {
		// Check for comments first.
		if cg, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range cg.List {
				if ok, args := extractArguments(comment.Text); ok {
					c.config.ParseConfig(args)
					lastComment = comment.Text // Save the comment if it contains "v:service"
				}
			}
		}

		// We're looking for type specifications (struct declarations).
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if the type is a struct.
		if st, ok := typeSpec.Type.(*ast.StructType); ok {
			// If there's a "d:service" comment, associate it with the struct.
			if lastComment != "" {
				// Reset the last comment after it is used.
				lastComment = ""

				holder := TypeHolder{
					typeSpec: typeSpec,
					embedded: []string{},
				}
				holder.embedded = c.findEmbeddedStructs(st, holder.embedded)

				structs = append(structs, holder)
			}
		}

		return true
	})

	return structs, nil
}

func (c *collect) findLocalStructType(structName string) *ast.StructType {
	var ret *ast.StructType

	// Look through all files for this struct
	err := filepath.Walk(c.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("error parsing file: %w", err)
		}

		ast.Inspect(astFile, func(n ast.Node) bool {
			// We're looking for type specifications (struct declarations).
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			// Check if the type is a struct.
			if st, ok := typeSpec.Type.(*ast.StructType); ok {
				if typeSpec.Name.Name == structName {
					ret = st
					return true
				}
			}
			return true
		})

		return nil
	})
	if err != nil {
		panic(err)
	}

	return ret
}

func (c *collect) findPackageStructType(pkg string, structName string) *ast.StructType {
	panic("implement findPackageStructType")
	return nil
}

// Recursive function to collect embedded structs from an ast.StructType
func (c *collect) findEmbeddedStructs(structType *ast.StructType, embeddedStructs []string) []string {
	for _, field := range structType.Fields.List {
		// Check if the field is an embedded struct (i.e., no field name)
		if len(field.Names) == 0 {
			// Get the type of the embedded struct
			switch t := field.Type.(type) {
			case *ast.Ident:
				// Simple embedded struct
				embeddedStructs = append(embeddedStructs, t.Name)
				s := c.findLocalStructType(t.Name)
				embeddedStructs = c.findEmbeddedStructs(s, embeddedStructs)
			case *ast.SelectorExpr:
				// For embedded structs from other packages (e.g., pkg.Struct)
				//pkg := t.X.(*ast.Ident).Name
				//name := t.Sel.Name
				//embeddedStructs = append(embeddedStructs, c.findPackageStructType(pkg, name))
			case *ast.StarExpr:
				// Dereferencing pointer types (e.g., *Struct or *pkg.Struct)
				switch expr := t.X.(type) {
				case *ast.Ident:
					embeddedStructs = append(embeddedStructs, expr.Name)
				case *ast.SelectorExpr:
					//pkg := expr.X.(*ast.Ident).Name
					//name := expr.Sel.Name
					//embeddedStructs = append(embeddedStructs, c.findPackageStructType(pkg, name))
				case *ast.StructType:
					// Recursively process embedded anonymous struct
					embeddedStructs = c.findEmbeddedStructs(expr, embeddedStructs)
				}
			case *ast.StructType:
				// Recursively process embedded anonymous struct
				embeddedStructs = c.findEmbeddedStructs(t, embeddedStructs)
			}
		}
	}
	return embeddedStructs
}
