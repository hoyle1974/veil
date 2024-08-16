package main

import (
	"context"

	"github.com/hoyle1974/veil/veil"
)

func VeilInitServer() {
	veil.RegisterService("main.Foo", func(s any, method string, args []any, reply *[]any) {
		if method == "Beep" {
			ret, err := s.(FooInterface).Beep(
				context.Background(),
				args[0].(int),
			)
			*reply = append(*reply, err)
			*reply = append(*reply, ret)
		}
	})
}
