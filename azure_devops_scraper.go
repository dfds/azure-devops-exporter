package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"time"
)

func channelBuildsResponseAsStringBetween(
	adoPersonalAccessToken string,
	projectIDs <-chan string,
	startTime time.Time,
	endTime time.Time) <-chan string {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	formattedStartTime := startTime.Format(time.RFC3339)
	formattedEndTime := endTime.Format(time.RFC3339)

	out := make(chan string)

	go func() {
		for projectID := range projectIDs {

			url := "https://dev.azure.com/dfds/" + projectID + "/_apis/build/builds?api-version=5.1&$top=5000&statusFilter=completed&minTime=" + formattedStartTime + "&maxTime=" + formattedEndTime
			resp, err := client.R().
				Get(url)

			panicOnError(err)
			out <- resp.String()
		}
		close(out)
	}()
	return out
}

func channelProjectIDs(adoPersonalAccessToken string) (<-chan string, int) {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	resp, err := client.R().
		Get("https://dev.azure.com/dfds/_apis/projects?api-version=5.1")

	panicOnError(err)

	projectsResponse := ProjectsResponse{}
	json.Unmarshal(resp.Body(), &projectsResponse)

	out := make(chan string)

	go func() {
		for i := 0; i < len(projectsResponse.Value); i++ {
			project := projectsResponse.Value[i]

			out <- project.ID
		}
		close(out)
	}()

	return out, len(projectsResponse.Value)
}

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
	} `json:"value"`
}
