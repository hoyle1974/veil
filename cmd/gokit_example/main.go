package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
)

type FooService struct {
}

func (f *FooService) SaySomething(ctx context.Context, name string, value int) (string, error) {
	ret := fmt.Sprintf("Hi %s, your value was %d", name, value)
	if value == -1 {
		return "", errors.New("you wanted an error")
	}
	return ret, nil
}

type FooService_Interface interface {
	SaySomething(ctx context.Context, name string, value int) (string, error)
}

type FooService_SaySomething_Request struct {
	Ctx   context.Context
	Name  string
	Value int
}

type FooService_SaySomething_Response struct {
	Arg1 string
	Err  error
}

func make_FooService_SaySomething_Endpoint(svc FooService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FooService_SaySomething_Request)
		arg1, err := svc.SaySomething(req.Ctx, req.Name, req.Value)
		return FooService_SaySomething_Response{arg1, err}, err
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeFooService_SaySomething_Request(_ context.Context, r *http.Request) (interface{}, error) {
	var request FooService_SaySomething_Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeFooService_SaySomething_Response(_ context.Context, r *http.Response) (FooService_SaySomething_Response, error) {
	var request FooService_SaySomething_Response
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "server" {
		server()
	}
	if os.Args[1] == "client" {
		client()
	}

	fmt.Println("Waiting.")
	select {}
}
