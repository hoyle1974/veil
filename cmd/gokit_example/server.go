package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoyle1974/veil/veil"
)

func server() {
	fmt.Println("-- server --")

	veil.VeilInitServer()

	// Make these visible remotely
	if err := veil.Serve(&BarService{}); err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":8181", nil))

}
