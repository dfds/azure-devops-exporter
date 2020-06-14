package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

func (diskStorage) storeScrapeResult(timeStamp time.Time, fileContent string) {

	dir, err := os.Getwd()
	check(err)

	fileName := "azure-devops-builds-" + timeStamp.Format(time.RFC3339) + ".json"

	filePathAndName := dir + "/scrape-results/" + fileName
	dataToWrite := []byte(fileContent)
	err = ioutil.WriteFile(filePathAndName, dataToWrite, 0644)
	check(err)
}

func (diskStorage) getLastScrapeStartTime() time.Time {

	dir, err := os.Getwd()
	check(err)
	root := dir + "/scrape-results/"

	fileInfo, err := ioutil.ReadDir(root)
	check(err)

	var scrapeTimes []string
	for _, file := range fileInfo {
		scrapeTime := strings.TrimLeft(strings.TrimRight(file.Name(), ".json"), "azure-devops-builds-")
		scrapeTimes = append(scrapeTimes, scrapeTime)
	}

	if len(scrapeTimes) == 0 {
		return time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	sort.Strings(scrapeTimes)

	last := scrapeTimes[len(scrapeTimes)-1]

	lastScrapeTime, _ := time.Parse("2006-01-02T15:04:05Z", last)

	return lastScrapeTime
}
