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

func GetDataForGoFile(fileName string, config Config) (VeilData, error) {
	fmt.Println("Getting Data for: ", fileName)

	data := VeilData{
		Filename:    fileName,
		PackageName: "main",
		Structs:     []Struct{},
		Types:       []string{},
	}

	// Let's do basic parse of the file, and find tags we care about
	taggedStructs, err := parseForTags(fileName, config)
	if err != nil {
		return data, err
	}
	config.TaggedStructs = taggedStructs

	projectDir := filepath.Dir(fileName)
	fmt.Println("Project Directory: ", projectDir)

	allStructs := map[string]Struct{}
	types := map[string]any{}

	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fmt.Println("Parsing: ", path)

		for _, structName := range taggedStructs {
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

			methods, tempTypes := getAllMethods(astFile, structName, types)
			types = tempTypes

			s.Methods = append(s.Methods, methods...)
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

func getAllMethods(astFile *ast.File, name string, types map[string]any) ([]Method, map[string]any) {
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

func parseForTags(fileName string, config Config) ([]string, error) {
	structs := []string{}

	astFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
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
					config.ParseConfig(args)
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
		if _, ok := typeSpec.Type.(*ast.StructType); ok {
			// If there's a "d:service" comment, associate it with the struct.
			if lastComment != "" {
				// Reset the last comment after it is used.
				structs = append(structs, typeSpec.Name.String())
				lastComment = ""
			}
		}

		return true
	})

	return structs, nil
}

/*
func findMethodsWithReceiverInProject(pkgName, fileName, structName string, config Config) (VeilData, error) {
	data := VeilData{}

	data.Filename = fileName
	data.PackageName = pkgName
	data.Structs = []Struct{}
	data.Types = []string{}

	// Assuming the project is in the current directory
	projectDir := filepath.Join("path", "projectName")

	var methods []*ast.FuncDecl
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		data, err = collectData2(pkgName, path, config, data)
		if err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error traversing project: %w", err)
	}

	return methods, nil
}
*/

func collectData(pkgName string, fileName string, config Config) (VeilData, error) {
	data := VeilData{}

	data.Filename = fileName
	data.PackageName = pkgName
	data.Structs = []Struct{}
	data.Types = []string{}

	types := map[string]any{}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return data, fmt.Errorf("error parsing directory: %w", err)
	}

	// Store the comments in the file.
	var lastComment string

	ast.Inspect(astFile, func(n ast.Node) bool {
		// Check for comments first.
		if cg, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range cg.List {
				if ok, args := extractArguments(comment.Text); ok {
					// if strings.Contains(comment.Text, "v:service") {
					// values := strings.Split(comment.Text, " ")
					config.ParseConfig(args)

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
		if _, ok := typeSpec.Type.(*ast.StructType); ok {
			// If there's a "d:service" comment, associate it with the struct.
			if lastComment != "" {
				// Reset the last comment after it is used.
				lastComment = ""

				// Generate method data structure
				methods := []Method{}
				for _, method := range GetMethodsForStruct(astFile, typeSpec.Name.Name) {
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

				data.Structs = append(data.Structs, Struct{
					Name:    typeSpec.Name.Name,
					Methods: methods,
				})

			}
		}

		return true
	})

	data.Types = extractTypes(types)

	return data, nil
}
