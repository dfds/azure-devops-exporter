package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	projectIDs := getProjectIDs(token)
	//projectIDs := []string{"136d92f4-a14a-422c-9f0e-230f6dbd90b1","785336a7-e841-46ba-b632-5092b88c7907"}

	storage := diskStorage{}

	buildStringsChannel := make(chan string)

	lastScrapeStartTime := storage.getLastScrapeStartTime()

	scrapeStartTime := time.Now().UTC()

	timeDifference := scrapeStartTime.Sub(lastScrapeStartTime)
	fmt.Println("Scraping finished builds between: '" + lastScrapeStartTime.Format(time.RFC3339) + "' and '" + scrapeStartTime.Format(time.RFC3339) + "' a HH:MM:SS " + time.Time{}.Add(timeDifference).Format("15:04:05") + " difference")

	// start a goroutine for each project
	for _, projectID := range projectIDs {
		go processProject(
			buildStringsChannel,
			token,
			projectID,
			lastScrapeStartTime,
			scrapeStartTime)
	}

	// Collect one projectBuildsStrings per project from the channel
	projectsBuildStrings := make([]string, 0)
	for i := 1; i <= len(projectIDs); i++ {
		projectBuildsString := <-buildStringsChannel
		if projectBuildsString != "{\"count\":0,\"value\":[]}" {
			projectBuildsString := removeWrapperObject(projectBuildsString)
			projectsBuildStrings = append(projectsBuildStrings, projectBuildsString)
		}
		fmt.Print(strconv.Itoa(i) + " ")
	}

	fileContent := "[" + strings.Join(projectsBuildStrings[:], ",") + "]"
	storage.storeScrapeResult(scrapeStartTime, fileContent)
}

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
	} `json:"value"`
}

type Storage interface {
	storeScrapeResult(timeStamp time.Time, fileContent string)
	getLastScrapeStarTime() time.Time
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processProject(
	buildStrings chan<- string,
	adoPersonalAccessToken string,
	projectID string,
	startTime time.Time,
	endTime time.Time) {

	buildsResponseAsString := getBuildsResponseAsStringBetween(
		adoPersonalAccessToken,
		projectID,
		startTime,
		endTime)

	buildStrings <- buildsResponseAsString
}

func removeWrapperObject(buildsResponseAsString string) string {

	i := strings.Index(buildsResponseAsString, "[")
	buildsResponseAsString = buildsResponseAsString[i+1:]

	buildsResponseAsString = strings.TrimSuffix(buildsResponseAsString, "]}")

	return buildsResponseAsString
}

func getBuildsResponseAsStringBetween(
	adoPersonalAccessToken string,
	projectID string,
	startTime time.Time,
	endTime time.Time) string {

	client := resty.New()
	// Bearer Auth Token for all request
	client.SetBasicAuth("", adoPersonalAccessToken)
	formattedStartTime := startTime.Format(time.RFC3339)
	formattedEndTime := endTime.Format(time.RFC3339)

	url := "https://dev.azure.com/dfds/" + projectID + "/_apis/build/builds?api-version=5.1&$top=5000&statusFilter=completed&minTime=" + formattedStartTime + "&maxTime=" + formattedEndTime
	resp, _ := client.R().
		Get(url)

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
