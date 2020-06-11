package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	projectIDs := getProjectIDs(token)

	fmt.Print(projectIDs)

	storage := diskStorage{}
	existingBuildsIDs := storage.getExistingBuildIDs()

	buildStringsChannel := make(chan string)

	// start a goroutine for each project
	for _, projectID := range projectIDs {
		go processProject(buildStringsChannel, storage, token, projectID, existingBuildsIDs)
	}

	// Collect one projectBuildsStrings per project from the channel
	projectsBuildStrings := make([]string, 0)
	for i := 1; i <= len(projectIDs); i++ {
		projectBuildsString := <-buildStringsChannel
		if projectBuildsString != "" {
			projectsBuildStrings = append(projectsBuildStrings, projectBuildsString)
		}
		fmt.Print(strconv.Itoa(i) + " ")
	}

	fileContent := "[" + strings.Join(projectsBuildStrings[:], ",") + "]"
	storage.storeScrapeResult(fileContent)
}

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
	} `json:"value"`
}

type Storage interface {
	getExistingBuildIDs() []string
	storeBuild(buildID string, fileContent string)
	storeScrapeResult(fileContent string)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processProject(
	buildStrings chan<- string,
	storage Storage,
	adoPersonalAccessToken string,
	projectID string,
	existingBuildIDs []string) {

	buildsResponseAsString := getBuildsResponseAsString(adoPersonalAccessToken, projectID)

	idToBuildMap := convertBuildsResponseToMap(buildsResponseAsString)

	idToBuildMap = removeExistingBuilds(existingBuildIDs, idToBuildMap)

	for buildID, javascriptObject := range idToBuildMap {
		storage.storeBuild(buildID, javascriptObject)
	}

	builds := make([]string, 0, len(idToBuildMap))

	for _, build := range idToBuildMap {
		builds = append(builds, build)
	}

	buildStrings <- strings.Join(builds[:], ",")
}

func removeExistingBuilds(existingBuildIDs []string, idToBuild map[string]string) map[string]string {
	for _, existingBuildID := range existingBuildIDs {
		delete(idToBuild, existingBuildID)

	}

	return idToBuild
}

func convertBuildsResponseToMap(buildsResponseAsString string) map[string]string {

	buildStrings := strings.Split(buildsResponseAsString, "{\"_links")
	buildStrings = buildStrings[1:]

	if len(buildStrings) == 0 {
		return nil
	}

	lastBuild := buildStrings[len(buildStrings)-1]
	buildStrings[len(buildStrings)-1] = lastBuild[:len(lastBuild)-2]

	regExPattern := regexp.MustCompile(`"id":(\d*),"buildNumber"`)

	var results = make(map[string]string)
	for _, currentString := range buildStrings {
		key := regExPattern.FindStringSubmatch(currentString)[1]

		results[key] = "{\"_links" + strings.TrimRight(currentString, ",")
	}

	return results
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
