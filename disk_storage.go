package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

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

func (diskStorage) storeScrapeResult(fileContent string) {

	dir, err := os.Getwd()
	check(err)

	t := time.Now().UTC()

	fileName := "azure-devops-builds-" + t.Format(time.RFC3339) + ".json"

	filePathAndName := dir + "/scrape-results/" + fileName
	dataToWrite := []byte(fileContent)
	err = ioutil.WriteFile(filePathAndName, dataToWrite, 0644)
	check(err)
}
