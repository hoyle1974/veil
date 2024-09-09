package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hoyle1974/veil/veil"
)

func testBarService(t *testing.T, bar BarService_Interface) {
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

func testBarServiceDown(t *testing.T, bar BarService_Interface) {
	tests := []struct {
		name     string
		value    int
		expected string
		err      string
	}{
		{"Jack", 123, "", "Post \"http://localhost:8181/BarService/SaySomething\": dial tcp [::1]:8181: connect: connection refused"},
		{"Bob", 5, "", "Post \"http://localhost:8181/BarService/SaySomething\": dial tcp [::1]:8181: connect: connection refused"},
		{"Jill", -1, "", "Post \"http://localhost:8181/BarService/SaySomething\": dial tcp [::1]:8181: connect: connection refused"},
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

func TestBarService(t *testing.T) {
	bar := &BarService{}

	testBarService(t, bar)
}

type MockConnection struct{}

func (c MockConnection) Get(path string, jsonData []byte) (*http.Response, error) {
	url := "http://localhost:8181" + path

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// Make the HTTP request
	return client.Do(req)
}

type MockConnFactory struct{}

func (m MockConnFactory) GetConnection() any {
	return MockConnection{}
}

type MockServFactory struct {
	mux *http.ServeMux
}

func (m MockServFactory) GetServer() any {
	return m.mux
}

func setupRemoteSuite(t testing.TB, enableNetwork bool) func(t testing.TB) {
	log.Println("setup suite")

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8181",
		Handler: mux,
	}

	control := veil.InitTestFramework(MockConnFactory{}, MockServFactory{mux: mux})
	control.StartTest(t)

	veil.Serve(&BarService{})

	if enableNetwork {
		go server.ListenAndServe()
	}

	// Return a function to teardown the test
	return func(t testing.TB) {
		control.StopTest(t)
		if enableNetwork {
			server.Shutdown(context.Background())
		}
	}
}

func TestBarService_Remote(t *testing.T) {
	teardownSuite := setupRemoteSuite(t, true)
	defer teardownSuite(t)

	bar, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		t.Errorf("Could not lookup BarService_Interface")
		return
	}
	if bar == nil {
		t.Errorf("Lookup BarService_Interface returned nil")
	}

	testBarService(t, bar)
}

func TestBarService_Remote_NetworkDown(t *testing.T) {
	teardownSuite := setupRemoteSuite(t, false)
	defer teardownSuite(t)

	bar, err := veil.Lookup[BarService_Interface]()
	if err != nil {
		t.Errorf("Could not lookup BarService_Interface")
		return
	}
	if bar == nil {
		t.Errorf("Lookup BarService_Interface returned nil")
	}

	testBarServiceDown(t, bar)
}
