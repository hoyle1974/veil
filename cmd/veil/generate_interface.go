package main

// GenerateInterfaceWithMethods generates a Go interface that includes all methods for a given struct type (from ast.TypeSpec).
// func GenerateInterfaceWithMethods(fqdn string, interfaceName string, file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
func (s *Source) GenerateInterface() error {

	// Generate interface name based on the struct name.
	var b = &s.file.common

	b.Sprintf("type %s interface {\n", s.InterfaceName())

	// Iterate over the methods and generate method signatures.
	for _, method := range GetMethodsForStruct(s.astFile, s.spec.Name.Name) {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {
			b.Sprintf("\t%s\n", methodSignature)
		}
	}

	b.Sprintf("}\n")

	return nil
}
