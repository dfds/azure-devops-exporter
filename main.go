package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
type awsStorage struct{}

func main() {

	// Get Azure Devops personal access token from environment
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")

	projectIDs := getProjectIDs(token)

	fmt.Print(projectIDs)

	storage := awsStorage{}
	existingBuildsIDs := storage.getExistingBuildIDs()
	var waitGroup sync.WaitGroup
	for _, projectID := range projectIDs {
		waitGroup.Add(1)
		go processProject(&waitGroup, storage, token, projectID, existingBuildsIDs)
	}

	waitGroup.Wait()

}

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

func (awsStorage) getExistingBuildIDs() []string {

	return nil
}

var myBucket = "dfds-datalake"

func (awsStorage) storeBuild(buildID string, fileContent string) {

	var awsConfig *aws.Config
	awsConfig = &aws.Config{
		Region: aws.String("eu-west-1"),
	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(awsConfig))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Create an uploader with the session and custom options
	// uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
	//     u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
	//     u.Concurrency = 2            // default is 5
	// })

	//defer f.Close()
	// Source sample https://paulbradley.org/gos3/
	buildFileKey := "azure-devops/" + buildID + ".json"
	dataToWrite := []byte(fileContent)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(buildFileKey),
		Body:   bytes.NewReader(dataToWrite),
	})

	//in case it fails to upload
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return
	}

	fmt.Printf("file uploaded to, %s\n", result.Location)
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

	for buildID, javascriptObject := range idToBuildMap {
		storage.storeBuild(buildID, javascriptObject)
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
