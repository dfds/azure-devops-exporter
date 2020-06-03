package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", token)
	resp, err := client.R().Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=5.1-preview.1")

	fmt.Println("Response Info:")
	fmt.Println("Error      :", err)
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Proto      :", resp.Proto())
	fmt.Println("Time       :", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())
	fmt.Println("Body       :\n", resp)
	fmt.Println()
}
