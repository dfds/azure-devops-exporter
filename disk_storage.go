package main

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type diskStorage struct{}

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
