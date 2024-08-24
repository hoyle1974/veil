package main

type File struct {
	head       Builder
	common     Builder
	clientInit Builder
	serverInit Builder
}

func NewFile(pkgName string) *File {
	f := &File{}
	f.head.Sprintf("package %s \n\n", pkgName)
	return f
}

func (f *File) String() string {
	var b Builder
	b.WriteString(f.head.String())
	b.WriteString(f.common.String())
	b.WriteString("func init() {\n")
	b.WriteString("veil.RegisterServerInit(func(){\n")
	b.WriteString(f.serverInit.String())
	b.WriteString("})\n")
	b.WriteString("veil.RegisterClientInit(func(){\n")
	b.WriteString(f.clientInit.String())
	b.WriteString("})\n")
	b.WriteString("}\n")

	return b.String()
}
