package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type AgentcloudsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID                            string `json:"id"`
		AgentCloudID                  int    `json:"agentCloudId"`
		Name                          string `json:"name"`
		Type                          string `json:"type"`
		Internal                      bool   `json:"internal"`
		AcquireAgentEndpoint          string `json:"acquireAgentEndpoint"`
		ReleaseAgentEndpoint          string `json:"releaseAgentEndpoint"`
		GetAgentDefinitionEndpoint    string `json:"getAgentDefinitionEndpoint"`
		GetAgentRequestStatusEndpoint string `json:"getAgentRequestStatusEndpoint"`
	} `json:"value"`
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", token)

	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=6.0-preview.1")

	agentcloudsResponse := AgentcloudsResponse{}
	json.Unmarshal(resp.Body(), &agentcloudsResponse)
	fmt.Printf("%#v\n", agentcloudsResponse)
}
