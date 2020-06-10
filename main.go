package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
)

type ProjectsResponse struct {
	Count int `json:"count"`
	Value []struct {
		ID string `json:"id"`
	} `json:"value"`
}

type Storage interface {
	getExistingBuildIDs() []string
	storeBuild(buildID string, fileContent string)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type diskStorage struct{}

func (diskStorage) getExistingBuildIDs() []string {
	dir, err := os.Getwd()
	check(err)
	root := dir + "/existing-builds/"

	fileInfo, err := ioutil.ReadDir(root)
	check(err)

	var existingBuilds []string
	for _, file := range fileInfo {
		existingBuilds = append(existingBuilds, strings.TrimRight(file.Name(), ".json"))
	}

	return existingBuilds
}
func (diskStorage) storeBuild(buildID string, fileContent string) {

	dir, err := os.Getwd()
	check(err)

	fmt.Print("storing: '" + buildID + "' ")

	fileName := dir + "/existing-builds/" + buildID + ".json"
	dataToWrite := []byte(fileContent)
	err = ioutil.WriteFile(fileName, dataToWrite, 0644)
	check(err)
}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	projectIDs := getProjectIDs(token)

	fmt.Print(projectIDs)

	storage := diskStorage{}
	existingBuildsIDs := storage.getExistingBuildIDs()
	var waitGroup sync.WaitGroup
	for _, projectID := range projectIDs {
		waitGroup.Add(1)
		go processProject(&waitGroup, storage, token, projectID, existingBuildsIDs)
	}

	waitGroup.Wait()
}

func processProject(
	waitGroup *sync.WaitGroup,
	storage Storage,
	adoPersonalAccessToken string,
	projectID string,
	existingBuildIDs []string) {

	defer waitGroup.Done()

	buildsResponseAsString := getBuildsResponseAsString(adoPersonalAccessToken, projectID)

	idToBuildMap := convertBuildsResponseToMap(buildsResponseAsString)

	idToBuildMap = removeExistingBuilds(existingBuildIDs, idToBuildMap)

	for buildId, javascriptObject := range idToBuildMap {
		storage.storeBuild(buildId, javascriptObject)
	}
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
