package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hoyle1974/veil/veil"
)

func TestBarService(t *testing.T) {
	bar := BarService{}

	tests := []struct {
		name     string
		value    int
		expected string
		err      string
	}{
		{"Jack", 123, "Hi Jack, your value was 123", ""},
		{"Bob", 5, "Hi Bob, your value was 5", ""},
		{"Jill", -1, "", "you wanted an error"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v,%v", test.name, test.value)
		t.Run(testname, func(t *testing.T) {
			ret, err := bar.SaySomething(context.Background(), test.name, test.value)

			errs := ""
			if err != nil {
				errs = err.Error()
			}
			if errs != test.err {
				t.Errorf("SaySomething returned error: [%v] vs [%v]", err, test.err)
				return
			}
			if ret != test.expected {
				t.Errorf("SaySomething returned the wrong value: %s vs %s", ret, test.expected)
				return
			}
		})
	}
}

func setupRemoteSuite(t testing.TB) func(t testing.TB) {
	log.Println("setup suite")

	control := veil.InitTestFramework(MockConnFactory{})
	control.StartTest(t)

	veil.Serve(&BarService{})

	server := &http.Server{Addr: ":8181"}
	go server.ListenAndServe()

	// Return a function to teardown the test
	return func(t testing.TB) {
		log.Println("teardown suite")
		control.StopTest(t)
		server.Shutdown(context.Background())
	}
}

type MockConnFactory struct{}

func (m MockConnFactory) GetConnection() any {
	return "http://localhost:8181"
}

func TestBarService_Remote(t *testing.T) {
	teardownSuite := setupRemoteSuite(t)
	defer teardownSuite(t)

	bar, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		t.Errorf("Could not lookup BarService_Interface")
		return
	}
	if bar == nil {
		t.Errorf("Lookup BarService_Interface returned nil")
	}

	tests := []struct {
		name     string
		value    int
		expected string
		err      string
	}{
		{"Jack", 123, "Hi Jack, your value was 123", ""},
		{"Bob", 5, "Hi Bob, your value was 5", ""},
		{"Jill", -1, "", "you wanted an error"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v,%v", test.name, test.value)
		t.Run(testname, func(t *testing.T) {
			ret, err := bar.SaySomething(context.Background(), test.name, test.value)

			errs := ""
			if err != nil {
				errs = err.Error()
			}
			if errs != test.err {
				t.Errorf("SaySomething returned error: [%v] vs [%v]", err, test.err)
				return
			}
			if ret != test.expected {
				t.Errorf("SaySomething returned the wrong value: %s vs %s", ret, test.expected)
				return
			}
		})
	}
}
