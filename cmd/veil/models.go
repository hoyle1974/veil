package main

type Arg struct {
	Name string // Name of the argument
	Type string // The type of the argument
}

// Describe a method
type Method struct {
	Name    string   // Name of the method
	Args    []Arg    // The arguments of the method
	Returns []string // The return value types
}

// Describes a struct that we want to expose publically
type Struct struct {
	Name    string   // The name of the struct
	Methods []Method // The methods that will be made public
}

type VeilData struct {
	Filename    string   // The original filename
	PackageName string   // The package
	Structs     []Struct // The structs that will be exposed remotely
	Packages    []string // Any packages that need to be included
}
