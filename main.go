package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type AgentcloudsRequestsResponse struct {
	Count int `json:"count"`
	Value []struct {
		AgentCloudID int    `json:"agentCloudId"`
		RequestID    string `json:"requestId"`
		Pool         struct {
			ID       int         `json:"id"`
			IsHosted bool        `json:"isHosted"`
			PoolType string      `json:"poolType"`
			Size     int         `json:"size"`
			IsLegacy interface{} `json:"isLegacy"`
			Options  interface{} `json:"options"`
		} `json:"pool"`
		Agent struct {
			ID                int         `json:"id"`
			Name              interface{} `json:"name"`
			Version           interface{} `json:"version"`
			Status            int         `json:"status"`
			ProvisioningState interface{} `json:"provisioningState"`
		} `json:"agent"`
		AgentSpecification struct {
			VMImage string `json:"VMImage"`
		} `json:"agentSpecification,omitempty"`
		AgentData struct {
			RequestID int `json:"RequestId"`
		} `json:"agentData"`
		ProvisionRequestTime time.Time   `json:"provisionRequestTime"`
		ProvisionedTime      interface{} `json:"provisionedTime"`
		AgentConnectedTime   time.Time   `json:"agentConnectedTime"`
		ReleaseRequestTime   time.Time   `json:"releaseRequestTime"`
	} `json:"value"`
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	agentcloudsResponse := getAgentcloudRequests(token)
	fmt.Printf("%#v\n", "Found "+strconv.Itoa(agentcloudsResponse.Count)+" agent cloud requests")
}

func getAgentcloudRequests(adoPersonalAccessToken string) AgentcloudsRequestsResponse {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)

	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds/1/requests?api-version=5.1-preview.1")

	agentcloudsResponse := AgentcloudsRequestsResponse{}
	json.Unmarshal(resp.Body(), &agentcloudsResponse)

	return agentcloudsResponse
}
