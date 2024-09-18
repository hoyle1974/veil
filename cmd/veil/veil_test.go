package main

import (
	"testing"
)

func TestBarService_Parse(t *testing.T) {

	// pkgName := "main"
	fileName := "../local_example/bar_service.go"
	config := Config{}

	data, err := GetDataForGoFile(fileName, config)
	if err != nil {
		t.Errorf("error collecting data: %v", err)
	}
	if len(data.Structs) != 1 {
		t.Error("not enough structs")
	}
	if data.Structs[0].Name != "BarService" {
		t.Error("BarService not found")
	}
	if len(data.Structs[0].Methods) != 4 {
		t.Errorf("not enough methods in BarService struct, only found %d", len(data.Structs[0].Methods))
	}
}

func TestBarService_Versioning(t *testing.T) {

	// pkgName := "main"
	fileName := "../versioning_example/bar_service.go"
	config := Config{}

	data, err := GetDataForGoFile(fileName, config)
	if err != nil {
		t.Errorf("error collecting data: %v", err)
	}
	if len(data.Structs) != 2 {
		t.Errorf("wrong number of structs, expected 2 but found %d", len(data.Structs))
	}
	if data.Structs[0].Name != "BarService" {
		t.Error("BarService not found")
	}
	if len(data.Structs[0].Methods) != 1 {
		t.Errorf("BarService method count should be 1, found %d", len(data.Structs[0].Methods))
	}
	if data.Structs[1].Name != "BarService2" {
		t.Error("BarService2 not found")
	}
	if len(data.Structs[1].Methods) != 1 {
		t.Errorf("BarService2 method count should be 1, found %d", len(data.Structs[1].Methods))
	}
}
