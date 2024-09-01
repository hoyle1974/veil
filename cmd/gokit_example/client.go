package main

import (
	"context"
	"fmt"

	"github.com/hoyle1974/veil/veil"
)

type ConnFactory struct {
}

func (c ConnFactory) GetConnection() any {
	return nil
}

func client() {
	fmt.Println("-- client --")

	veil.VeilInitClient(ConnFactory{})

	bar, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		panic(err)
	}

	result, err := bar.SaySomething(context.Background(), "Jack", 2431)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", result)
}
