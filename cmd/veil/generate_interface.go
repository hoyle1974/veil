package main

import (
	"go/ast"
)

func (s *Source) GenerateInterface() error {
	data, err := GenerateInterfaceWithMethods(s.FQDN(), s.InterfaceName(), s.astFile, s.spec)
	if err != nil {
		return err
	}

	var b = &s.file.common

	b.WriteString("// Generated from " + s.fileName + "\n")
	b.WriteString(data)
	b.WriteString("\n")

	return nil
}

// GenerateInterfaceWithMethods generates a Go interface that includes all methods for a given struct type (from ast.TypeSpec).
func GenerateInterfaceWithMethods(fqdn string, interfaceName string, file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Generate interface name based on the struct name.

	var b Builder

	b.Sprintf("type %s interface {\n", interfaceName)

	// Iterate over the methods and generate method signatures.
	for _, method := range GetMethodsForStruct(file, typeSpec.Name.Name) {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {
			b.Sprintf("\t%s\n", methodSignature)
		}
	}

	b.Sprintf("}\n")

	return b.String(), nil
}
