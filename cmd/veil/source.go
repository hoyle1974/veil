package main

import (
	"fmt"
	"go/ast"
)

type Source struct {
	pkgName  string
	spec     *ast.TypeSpec
	astFile  *ast.File
	fileName string
	file     *File
}

func NewSource(pkgName string, spec *ast.TypeSpec, astFile *ast.File, fileName string, file *File) (*Source, error) {
	_, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("%s is not a struct", spec.Name.Name)
	}

	return &Source{pkgName: pkgName, spec: spec, astFile: astFile, fileName: fileName, file: file}, nil

}

func (s *Source) InterfaceName() string {
	return fmt.Sprintf("%sInterface", s.spec.Name.Name)
}

func (s *Source) RemoteName() string {
	return fmt.Sprintf("%sRemoteImpl", s.spec.Name.Name)
}

func (s *Source) FQDN() string {
	return s.pkgName + "." + s.spec.Name.Name
}

func (s *Source) Generate() {
	s.file.head.WriteString("import (\n")
	for _, i := range getImports(s.astFile) {
		s.file.head.Sprintf("	\"%s\"\n", i)
	}
	s.file.head.WriteString(")\n\n")

	s.file.serverInit.Sprintf("veil.RegisterRemoteImpl(&%s{})\n", s.RemoteName())

	s.GenerateInterface()
	s.GenerateRemoteImpl()
	s.GenerateServiceBindings()

}
