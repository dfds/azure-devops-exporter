package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	token := os.Getenv("ADO_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		panic("A valid Azure DevOps access token needs to be set in the environment variable: 'ADO_PERSONAL_ACCESS_TOKEN'")
	}
	//projectIDs := []string{"136d92f4-a14a-422c-9f0e-230f6dbd90b1","785336a7-e841-46ba-b632-5092b88c7907"}

	storage := awsStorage{}
	//storage := diskStorage{}

	lastScrapeStartTime := storage.getLastScrapeStartTime()
	//lastScrapeStartTime, _ := time.Parse("2006-01-02T15:04:05Z", "2020-06-18T06:45:19Z")

	scrapeStartTime := time.Now().UTC()
	//scrapeStartTime, _ := time.Parse("2006-01-02T15:04:05Z", " 2020-06-18T06:50:30Z")

	timeDifference := scrapeStartTime.Sub(lastScrapeStartTime)
	fmt.Println("Scraping finished builds between: '" + lastScrapeStartTime.Format(time.RFC3339) + "' and '" + scrapeStartTime.Format(time.RFC3339) + "' a HH:MM:SS " + time.Time{}.Add(timeDifference).Format("15:04:05") + " difference")

	// start a goroutine for each project

	projectIDsChannel, projectIDsCount := channelProjectIDs(token)

	fmt.Println(fmt.Sprintf("will process %d projects\n", projectIDsCount))

	buildsResponseAsStringsChannel := channelBuildsResponseAsStringBetween(
		token,
		projectIDsChannel,
		lastScrapeStartTime,
		scrapeStartTime)

	printProgressBar(projectIDsCount)

	// Collect one projectBuildsStrings per project from the channel
	projectsBuildStrings := make([]string, 0)
	for i := 1; i <= projectIDsCount; i++ {
		projectBuildsString := <-buildsResponseAsStringsChannel

		if projectBuildsString != "{\"count\":0,\"value\":[]}" {
			projectBuildsString := removeWrapperObject(projectBuildsString)
			projectsBuildStrings = append(projectsBuildStrings, projectBuildsString)
		}
		fmt.Print("█")
	}
	fmt.Println()

	if len(projectsBuildStrings) == 0 {
		return
	}

	fileContent := "[" + strings.Join(projectsBuildStrings[:], ",") + "]"
	storage.storeScrapeResult(scrapeStartTime, fileContent)
	fmt.Println()
}

func printProgressBar(length int) {
	fmt.Println("Progress:")

	progressBar := "┌"

	for i := 3; i <= length; i++ {
		progressBar = progressBar + "─"
	}
	progressBar = progressBar + "┐"
	rune := []rune(progressBar)
	rune[len(rune)/2] = '┬'
	progressBar = string(rune)
	fmt.Println(progressBar)
}

type Storage interface {
	storeScrapeResult(timeStamp time.Time, fileContent string)
	getLastScrapeStarTime() time.Time
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func removeWrapperObject(buildsResponseAsString string) string {

	i := strings.Index(buildsResponseAsString, "[")
	buildsResponseAsString = buildsResponseAsString[i+1:]

	buildsResponseAsString = strings.TrimSuffix(buildsResponseAsString, "]}")

	return buildsResponseAsString
}
