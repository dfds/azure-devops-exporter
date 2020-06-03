package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type agentcloud struct {
	agentCloudId                  int    `json:"agentCloudId"`
	internal                      bool   `json:"internal"`
	name                          string `json:"name"`
	acquireAgentEndpoint          string `json:"acquireAgentEndpoint"`
	releaseAgentEndpoint          string `json:"releaseAgentEndpoint"`
	getAgentDefinitionEndpoint    string `json:"getAgentDefinitionEndpoint"`
	getAgentRequestStatusEndpoint string `json:"getAgentRequestStatusEndpoint"`
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", token)
	// resp, err := client.R().Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=6.0-preview.1")

	//fmt.Print(resp.Result())
	client.R().SetResult(agentcloud{})
	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=6.0-preview.1")
	fmt.Printf("%#v\n", resp.Result())

	// fmt.Println("Response Info:")
	// fmt.Println("Error      :", err)
	// fmt.Println("Status Code:", resp.StatusCode())
	// fmt.Println("Status     :", resp.Status())
	// fmt.Println("Proto      :", resp.Proto())
	// fmt.Println("Time       :", resp.Time())
	// fmt.Println("Received At:", resp.ReceivedAt())
	// fmt.Println("Body       :\n", resp)
	// fmt.Println()
}
