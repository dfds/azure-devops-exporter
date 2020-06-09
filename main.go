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

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
		// Name           string    `json:"name"`
		// Description    string    `json:"description,omitempty"`
		// URL            string    `json:"url"`
		// State          string    `json:"state"`
		// Revision       int       `json:"revision"`
		// Visibility     string    `json:"visibility"`
		// LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"value"`
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	//agentcloudsResponse := getAgentcloudRequests(token)
	//fmt.Printf("%#v\n", "Found "+strconv.Itoa(agentcloudsResponse.Count)+" agent cloud requests")

	// jobRequests := getJobRequests(token, 169)

	// fmt.Printf("%#v\n", "Found "+strconv.Itoa(jobRequests.Count)+" agent requests")

	// buildStatistics := ConvertAgentRequestsResponseToBuildStatistics(jobRequests)

	// fmt.Print(buildStatistics)

	projectIDs := getProjectIDs(token)

	fmt.Print(projectIDs)
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

// This one
func ConvertAgentRequestsResponseToBuildStatistics(jobRequestsResponse JobRequestsResponse) []BuildStatistic {
	var buildStatistics []BuildStatistic

	for _, jobRequest := range jobRequestsResponse.Value {

		buildStatistic := BuildStatistic{
			StartTime:   jobRequest.QueueTime,
			EndTime:     jobRequest.FinishTime,
			QueueLength: jobRequest.ReceiveTime.Sub(jobRequest.QueueTime),
			BuildLength: jobRequest.FinishTime.Sub(jobRequest.ReceiveTime),
		}

		buildStatistics = append(buildStatistics, buildStatistic)
	}
	return buildStatistics
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

func getJobRequests(adoPersonalAccessToken string, poolId int) JobRequestsResponse {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/distributedtask/pools/" + strconv.Itoa(poolId) + "/jobrequests")

	jobRequestsResponse := JobRequestsResponse{}
	json.Unmarshal(resp.Body(), &jobRequestsResponse)

	return jobRequestsResponse
}

func getProjectIDs(adoPersonalAccessToken string) []string {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/_apis/projects?api-version=5.1")

	projectsResponse := ProjectsResponse{}
	json.Unmarshal(resp.Body(), &projectsResponse)

	projectIDCollection := []string{}
	for i := 0; i < len(projectsResponse.Value); i++ {
		project := projectsResponse.Value[i]

		projectIDCollection = append(projectIDCollection, project.ID)

	}
	return projectIDCollection
}

type JobRequestsResponse struct {
	Count int `json:"count"`
	Value []struct {
		RequestID     int       `json:"requestId"`
		QueueTime     time.Time `json:"queueTime"` // When the pool is assigned the job and waiting in queue
		AssignTime    time.Time `json:"assignTime"`
		ReceiveTime   time.Time `json:"receiveTime"` // Agent acknowledges that it has received the job
		FinishTime    time.Time `json:"finishTime"`  // Agent is finished with the request
		Result        string    `json:"result"`
		ServiceOwner  string    `json:"serviceOwner"`
		HostID        string    `json:"hostId"`
		ScopeID       string    `json:"scopeId"`
		PlanType      string    `json:"planType"`
		PlanID        string    `json:"planId"`
		JobID         string    `json:"jobId"`
		Demands       []string  `json:"demands"`
		ReservedAgent struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				Web struct {
					Href string `json:"href"`
				} `json:"web"`
			} `json:"_links"`
			ID                int    `json:"id"`
			Name              string `json:"name"`
			Version           string `json:"version"`
			OsDescription     string `json:"osDescription"`
			Enabled           bool   `json:"enabled"`
			Status            string `json:"status"`
			ProvisioningState string `json:"provisioningState"`
			AccessPoint       string `json:"accessPoint"`
		} `json:"reservedAgent"`
		Definition struct {
			Links struct {
				Web struct {
					Href string `json:"href"`
				} `json:"web"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"definition"`
		Owner struct {
			Links struct {
				Web struct {
					Href string `json:"href"`
				} `json:"web"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"owner"`
		Data struct {
			ParallelismTag string `json:"ParallelismTag"`
			IsScheduledKey string `json:"IsScheduledKey"`
		} `json:"data,omitempty"`
		PoolID             int           `json:"poolId"`
		AgentDelays        []interface{} `json:"agentDelays"`
		AgentSpecification struct {
			VMImage string `json:"vmImage"`
		} `json:"agentSpecification,omitempty"`
		OrchestrationID        string `json:"orchestrationId"`
		MatchesAllAgentsInPool bool   `json:"matchesAllAgentsInPool"`
	} `json:"value"`
}
