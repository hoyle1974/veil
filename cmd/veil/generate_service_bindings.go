package main

import (
	"fmt"
	"go/ast"
	"strings"
)

/*
func VeilInitServer() {
	veil.RegisterService("main.Foo", func(s any, method string, args []any, reply *[]any) {
		if method == "Beep" {
			ret, err := s.(FooInterface).Beep(
				context.Background(),
				args[0].(int),
			)
			*reply = append(*reply, err)
			*reply = append(*reply, ret)
		}
		if method == "Boop" {
			ret, err := s.(FooInterface).Boop(
				context.Background(),
				args[0].(string),
			)
			*reply = append(*reply, err)
			*reply = append(*reply, ret)
		}
	})
	go veil.StartServices()
}
*/

func GenerateServiceBindings(file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Ensure the type is a struct.
	_, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return "", fmt.Errorf("%s is not a struct", typeSpec.Name.Name)
	}

	// Generate interface name based on the struct name.
	interfaceName := fmt.Sprintf("%sInterface", typeSpec.Name.Name)
	var builder strings.Builder

	builder.WriteString("func VeilInitServer() {\n")
	builder.WriteString("	veil.RegisterService(\"main.Foo\", func(s any, method string, args []any, reply *[]any) {")

	methods := GetMethodsForStruct(file, typeSpec.Name.Name)

	for _, method := range methods {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {

			builder.WriteString("if method == \"" + method.Name.Name + "\" {\n")
			builder.WriteString("	ret, err := s.(" + interfaceName + ")." + method.Name.Name + "(\n")
			builder.WriteString("		context.Background(),\n")
			// builder.WriteString("		args[0].(int),\n")
			for i, param := range method.Type.Params.List {
				tas := getTypeAsString(param.Type)
				if i != 0 {
					builder.WriteString("		args[" + fmt.Sprintf("%d", i-1) + "].(" + tas + "),\n")
				}
			}
			builder.WriteString("	)\n")
			builder.WriteString("	*reply = append(*reply, ret)\n")
			builder.WriteString("	*reply = append(*reply, err)\n")
			builder.WriteString("}\n")
		}
	}

	builder.WriteString("	})\n")
	builder.WriteString("	go veil.StartServices()\n")
	builder.WriteString("}\n")

	return builder.String(), nil
}
