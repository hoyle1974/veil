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
