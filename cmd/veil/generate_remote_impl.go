package main

import (
	"fmt"
	"go/ast"
)

func (s *Source) GenerateRemoteImpl() error {
	data, err := GenerateRemoteImpl(s.FQDN(), s.RemoteName(), s.astFile, s.spec)
	if err != nil {
		return err
	}

	var b = &s.file.common
	b.Sprintf("// Generated from %s\n", s.fileName)
	b.WriteString(data)
	b.WriteString("\n")
	return nil
}

/*

	request := veil.Request{
		Service: "main.Foo",
		Method:  "Boop",
		Args:    []any{value},
	}

	reply := []any{}
	var result0 error
	var result1 string

	err := veil.Call(request, &reply)
	if err != nil {
		result0 = err
	} else {
		result0 = veil.NilGet[error](reply[0])
		result1 = veil.NilGet[string](reply[1])
	}

	return result1, result0

*/

func generateRequestObject(typeSpec *ast.TypeSpec, method *ast.FuncDecl) (string, string) {
	var b Builder

	// Create requests object
	reqObjName := fmt.Sprintf("%s_%s_Request", typeSpec.Name.Name, method.Name)
	b.Sprintf("type %s struct {\n", reqObjName)
	for i, param := range method.Type.Params.List {
		for _, name := range param.Names {
			if i != 0 {
				tas := getTypeAsString(param.Type)
				tag := "`bson:\"" + name.Name + "\"`"

				b.Sprintf("%s %s %s\n", UppercaseFirst(name.Name), tas, tag)
			}
		}
	}
	b.WriteString("}\n")

	return b.String(), reqObjName
}

func GenerateRemoteImpl(fqdn string, remoteImplName string, file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Generate interface name based on the struct name.

	var b Builder

	b.Sprintf("type %s struct {}\n", remoteImplName)

	// Iterate over the methods and generate method signatures.
	for _, method := range GetMethodsForStruct(file, typeSpec.Name.Name) {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {

			requestDecl, requestObjName := generateRequestObject(typeSpec, method)
			b.WriteString(requestDecl)

			temp := ""
			for i, param := range method.Type.Params.List {
				for _, name := range param.Names {
					if i != 0 {
						if i > 1 {
							temp += ","
						}
						temp += name.Name
					}
				}
			}

			b.Sprintf("func (r *%s) %s {\n", remoteImplName, methodSignature)
			b.Sprintf("data, err := bson.Marshal(%s{ %s})\n", requestObjName, temp)
			b.Sprintf("if err != nil {\n")
			b.Sprintf("	panic(err)\n")
			b.Sprintf("}\n")
			b.Sprintf("request := veil.Request{\n")
			b.Sprintf("	Service: \"%s\",\n", fqdn)
			b.Sprintf("	Method:  \"%s\",\n", method.Name.String())
			b.Sprintf("	Args:    data,\n")
			b.Sprintf("}\n")
			b.Sprintf("reply := []any{}\n")

			// Handle the return values.
			last := ""
			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					b.Sprintf("var result%d %s\n", i, tas)
					last = fmt.Sprintf("result%d", i)
				}
			}
			b.Sprintf("err = veil.Call(request, &reply)\n")

			b.Sprintf("if err != nil {\n")
			b.Sprintf("	%s = err\n", last)
			b.Sprintf("} else {\n")
			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					b.Sprintf("result%d = veil.NilGet[%s](reply[%d])\n", i, tas, i)
				}
			}
			b.Sprintf("}\n")

			b.Sprintf("return ")
			if method.Type.Results != nil {
				for i := range method.Type.Results.List {
					if i == 0 {
						b.Sprintf("result%d", i)
					} else {
						b.Sprintf(", result%d", i)

					}
				}
			}
			b.Sprintf("\n")

			b.Sprintf("}\n")
		}
	}

	return b.String(), nil
}
