package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func client() {

	url := "http://localhost:8181/BarService/SaySomething"

	// Marshal the JSON data into a byte buffer
	jsonData, err := json.Marshal(BarService_SaySomething_Request{Name: "Jack", Value: 23})
	if err != nil {
		fmt.Println("error marshalling JSON:", err)
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error creating request:", err)
		return

	}

	// Set the content type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client
	client := &http.Client{}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code", resp.StatusCode)
		return
	}

	req2, _ := decodeBarService_SaySomething_Response(context.Background(), resp)
	fmt.Println(req2)

	// Read the response body
	defer resp.Body.Close()
	// ... (process response body if needed)
	fmt.Println("Request successful!")
}
