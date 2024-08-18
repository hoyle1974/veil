package main

import (
	"fmt"
	"go/ast"
	"strings"
)

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

func GenerateRemoteImpl(fqdn string, file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Ensure the type is a struct.
	_, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return "", fmt.Errorf("%s is not a struct", typeSpec.Name.Name)
	}

	// Generate interface name based on the struct name.
	implName := fmt.Sprintf("%sRemoteImpl", typeSpec.Name.Name)
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("type %s struct {}\n", implName))

	// Retrieve all methods associated with the struct.
	methods := GetMethodsForStruct(file, typeSpec.Name.Name)

	// Iterate over the methods and generate method signatures.
	for _, method := range methods {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {

			// Create requests object
			reqObjName := fmt.Sprintf("%s_%s_Request", typeSpec.Name.Name, method.Name)
			builder.WriteString("type " + reqObjName + " struct {\n")
			for i, param := range method.Type.Params.List {
				for _, name := range param.Names {
					if i != 0 {
						tas := getTypeAsString(param.Type)
						tag := "`bson:\"_" + name.Name + "\"`"

						builder.WriteString("D" + name.Name + " " + tas + " " + tag + "\n")
					}
				}
			}
			builder.WriteString("}\n")

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

			builder.WriteString(fmt.Sprintf("func (r *%s) %s {\n", implName, methodSignature))
			builder.WriteString("data, err := bson.Marshal(" + reqObjName + "{ " + temp + "})\n")
			builder.WriteString("if err != nil {\n")
			builder.WriteString("	panic(err)\n")
			builder.WriteString("}\n")
			builder.WriteString("request := veil.Request{\n")
			builder.WriteString("	Service: \"" + fqdn + "\",\n")
			builder.WriteString("	Method:  \"" + method.Name.String() + "\",\n")
			builder.WriteString("	Args:    data,\n")
			builder.WriteString("}\n")
			builder.WriteString("reply := []any{}\n")

			// Handle the return values.
			last := ""
			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					builder.WriteString(fmt.Sprintf("var result%d %s\n", i, tas))
					last = fmt.Sprintf("result%d", i)
				}
			}
			builder.WriteString("err = veil.Call(request, &reply)\n")

			builder.WriteString("if err != nil {\n")
			builder.WriteString("	" + last + " = err\n")
			builder.WriteString("} else {\n")
			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					builder.WriteString(fmt.Sprintf("result%d = veil.NilGet[%s](reply[%d])\n", i, tas, i))
				}
			}
			builder.WriteString("}\n")

			builder.WriteString("return ")
			if method.Type.Results != nil {
				for i := range method.Type.Results.List {
					if i == 0 {
						builder.WriteString(fmt.Sprintf("result%d", i))
					} else {
						builder.WriteString(fmt.Sprintf(", result%d", i))

					}
				}
			}

			builder.WriteString("}\n")
		}
	}

	return builder.String(), nil
}
