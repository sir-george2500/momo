package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

func CreateApiUser() (string, int, string, error) {
	url := "https://sandbox.momodeveloper.mtn.com/v1_0/apiuser"
	method := "POST"

	// Define the payload
	payload := strings.NewReader(`{
		"providerCallbackHost": "callbacks-do-not-work-in-sandbox.com"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Generate and set the X-Reference-Id
	xReferenceID := uuid.New().String()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Key", "c91b4295001440efb7d286f452575de7")
	req.Header.Add("X-Reference-Id", xReferenceID)

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", res.StatusCode, "", fmt.Errorf("failed to read response: %w", err)
	}

	return xReferenceID, res.StatusCode, string(body), nil
}

func CreateApiKey(xReferenceID string) (int, string, error) {
	url := "https://sandbox.momodeveloper.mtn.com/v1_0/apiuser/" + xReferenceID + "/apikey"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Key", "c91b4295001440efb7d286f452575de7")

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, "", fmt.Errorf("failed to read response: %w", err)
	}

	return res.StatusCode, string(body), nil
}

func CreateAccessToken() (int, string, error) {
	url := "https://sandbox.momodeveloper.mtn.com/collection/token/"
	method := "POST"

	// Use the exact same credentials that work in curl
	apiUserId := "0511e7d5-175c-4523-a3ff-7039e2c33233"
	apiKey := "52ebed8b483b4d8cb6ef904a6f05f1a9"

	// Set the payload
	payload := strings.NewReader("grant_type=client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Ocp-Apim-Subscription-Key", "c91b4295001440efb7d286f452575de7")

	// Set Basic Auth
	req.SetBasicAuth(apiUserId, apiKey)

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	// Read and return response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, "", fmt.Errorf("failed to read response: %w", err)
	}

	return res.StatusCode, string(body), nil
}

// func main() {
// 	http.HandleFunc("/call-api", handleRequest)
// 	http.HandleFunc("/health", healthCheckHandler) // Register health check endpoint
//
// 	fmt.Println("Server is running on port 8080...")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		fmt.Println("Failed to start server:", err)
// 	}
// }

func main() {
	// Step 1: Create API User
	xReferenceID, statusCode, responseBody, err := CreateApiUser()
	if err != nil {
		fmt.Printf("Error creating API user: %v\n", err)
		return
	}

	fmt.Printf("API User Created\nX-Reference-Id: %s\nResponse: %s\n", xReferenceID, responseBody)

	// Step 2: Create API Key
	if statusCode == http.StatusCreated {
		keyStatus, keyResponse, err := CreateApiKey(xReferenceID)
		if err != nil {
			fmt.Printf("Error creating API key: %v\n", err)
			return
		}

		fmt.Printf("API Key Created\nStatus: %d\nResponse: %s\n", keyStatus, keyResponse)

		fmt.Printf("The key was generated or created\n")

		// Step 3: Create Access Token using the API Key
		if keyStatus == 201 {
			// Extract the API key from the keyResponse (you would need to parse the response body here)
			// Assuming `keyResponse` contains the API key, for now, we can just use it as a string.
			// You might need to parse JSON response to extract the actual key.
			fmt.Print("moving to create the CreateAccessToken")
			// Here we just pass it directly assuming it’s in the response.
			// Adjust this to extract the actual API key from the response body
			// Step 4: Create Access Token
			tokenStatus, tokenResponse, err := CreateAccessToken()
			if err != nil {
				fmt.Printf("Error creating access token: %v\n", err)
				return
			}

			fmt.Printf("Access Token Created\nStatus: %d\nResponse: %s\n", tokenStatus, tokenResponse)
		}
	} else {
		fmt.Printf("API user creation failed with status: %d\n", statusCode)
	}
}
