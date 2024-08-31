package main

import (
	"log"
	"net/http"

	"github.com/hoyle1974/veil/veil"
)

func server() {
	// svc := FooService{}

	// FooService_SaySomething_Handler := httptransport.NewServer(
	// 	make_FooService_SaySomething_Endpoint(svc),
	// 	decodeFooService_SaySomething_Request,
	// 	encodeResponse,
	// )

	// http.Handle("/FooService/SaySomething", FooService_SaySomething_Handler)
	// log.Fatal(http.ListenAndServe(":8181", nil))

	veil.VeilInitServer()

	// Make these visible remotely
	if err := veil.Serve(&BarService{}); err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":8181", nil))

}
