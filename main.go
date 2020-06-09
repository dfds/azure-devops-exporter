package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
	} `json:"value"`
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	projectIDs := getProjectIDs(token)

	fmt.Print(projectIDs)

	r := getBuildsResponseAsString(token, "b37a3f59-1490-4321-b288-a7985e3da04c")
	fmt.Print(r)

}
func getBuildsResponseAsString(adoPersonalAccessToken string, projectID string) string {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	resp, _ := client.R().
		Get("https://dev.azure.com/dfds/" + projectID + "/_apis/build/builds?api-version=5.1&$top=5000")

	return resp.String()
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
