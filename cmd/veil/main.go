package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed rpc_service.tmpl
var rpc_service []byte

//go:embed gokit_service.tmpl
var gokit_service []byte

func getImports(file *ast.File) []string {
	var importPaths []string
	for _, imp := range file.Imports {
		importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
		importPaths = append(importPaths, importPath)
	}
	return importPaths
}

func title(in string) string {
	return cases.Title(language.English).String(in)
}
func lastItemIndex(a any) int {
	items := a.([]string)
	return len(items) - 1
}

type Arg struct {
	Name string
	Type string
}

type Method struct {
	Name    string
	Args    []Arg
	Returns []string
}

type Struct struct {
	Name           string
	RPCName        string
	InterfaceName  string
	RemoteImplName string
	Methods        []Method
}

type Data struct {
	Filename    string
	PackageName string
	Structs     []Struct
	Packages    []string
	Name        string
}

func extractArguments(input string) (bool, string) {
	// Check if "v:service" is present
	if !strings.Contains(input, "v:service") {
		return false, ""
	}

	// Split the string after "v:service"
	parts := strings.Split(input, "v:service ")
	if len(parts) < 2 {
		return false, ""
	}

	// Extract the arguments
	//arguments := strings.Fields(parts[1])

	return true, parts[1] //arguments
}

func main() {
	// templateFile := "rpc_service.tmpl"
	templateFile := "gokit_service.tmpl"

	tmpl := template.New(templateFile)
	tmpl.Funcs(map[string]any{
		"title":         title,
		"lastItemIndex": lastItemIndex,
	})

	data := Data{}

	// Replace "your/project/path" with the actual path to your project
	// projectPath :=  "/Users/jstrohm/code/veil/cmd/veil"
	// fmt.Println(os.Environ())
	fileName := os.Getenv("GOFILE")
	pkgName := os.Getenv("GOPACKAGE")

	if fileName == "" {
		fileName = "/Users/jstrohm/code/veil/cmd/gokit_example/bar_service.go"
		pkgName = "main"
	}

	data.Filename = fileName
	data.PackageName = pkgName
	data.Structs = []Struct{}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing directory:", err)
		return
	}

	// Store the comments in the file.
	var lastComment string

	config := lookupConfig()

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
					Name:           typeSpec.Name.Name,
					InterfaceName:  fmt.Sprintf("%s_Interface", typeSpec.Name.Name),
					RemoteImplName: fmt.Sprintf("%s_RemoteImpl", typeSpec.Name.Name),
					RPCName:        fmt.Sprintf("%s_RPC", typeSpec.Name.Name),
					Methods:        methods,
				})

			}
		}

		return true
	})

	tmpl, err = tmpl.Parse(config.GetTemplateString())
	if err != nil {
		panic(err)
	}

	ifile := config.Directory + "/" + "impl_" + fileName

	f, err := os.OpenFile(ifile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("goimports", "-w", ifile)
	err = cmd.Run()
	if err != nil {
		panic(fmt.Errorf("can't execute goimports on %s: %w", ifile, err))
	}

}
