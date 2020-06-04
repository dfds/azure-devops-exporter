package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type BuildStatistic struct {
	VirtualMachineImage string
	StartTime           time.Time
	EndTime             time.Time
	BuildLength         time.Duration
	QueueLength         time.Duration
}

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
	//https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/requests/list?view=azure-devops-rest-6.0
	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/distributedtask/agentclouds/1/requests?api-version=6.0-preview.1")

	agentcloudsResponse := AgentcloudsRequestsResponse{}
	json.Unmarshal(resp.Body(), &agentcloudsResponse)

	return agentcloudsResponse
}

func ConvertAgentcloudsRequestsResponseToBuildStatistics(agentcloudsRequestsResponse AgentcloudsRequestsResponse) []BuildStatistic {
	var buildStatistics []BuildStatistic

	for _, agentcloudsRequest := range agentcloudsRequestsResponse.Value {

		buildStatistic := BuildStatistic{
			StartTime:   agentcloudsRequest.ProvisionRequestTime,
			QueueLength: agentcloudsRequest.AgentConnectedTime.Sub(agentcloudsRequest.ProvisionRequestTime),
			BuildLength: agentcloudsRequest.ReleaseRequestTime.Sub(agentcloudsRequest.AgentConnectedTime),
		}

		buildStatistics = append(buildStatistics, buildStatistic)
	}
	return buildStatistics
}
