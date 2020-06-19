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

	//storage := awsStorage{}
	storage := diskStorage{}

	lastScrapeStartTime := storage.getLastScrapeStartTime()
	//lastScrapeStartTime, _ := time.Parse("2006-01-02T15:04:05Z", "2020-06-18T06:45:19Z")

	scrapeStartTime := time.Now().UTC()
	//scrapeStartTime, _ := time.Parse("2006-01-02T15:04:05Z", " 2020-06-18T06:50:30Z")

	timeDifference := scrapeStartTime.Sub(lastScrapeStartTime)
	fmt.Println("Scraping finished builds between: '" + lastScrapeStartTime.Format(time.RFC3339) + "' and '" + scrapeStartTime.Format(time.RFC3339) + "' a HH:MM:SS " + time.Time{}.Add(timeDifference).Format("15:04:05") + " difference")

	// start a goroutine for each project

	projectIDsChannel, projectIDsCount := channelProjectIDs(token)
	<-projectIDsChannel
	fmt.Println(fmt.Sprintf("Will process %d projects\n", projectIDsCount))
	printProgressBar(projectIDsCount)
	buildsResponseAsStringsChannel := channelBuildsResponseAsStringBetween(
		token,
		projectIDsChannel,
		lastScrapeStartTime,
		scrapeStartTime)

	buildsResponseWithoutEmptyChannel := channelRemoveEmptyResults(buildsResponseAsStringsChannel)
	buildsWithNoWrappersChannel := channelRemoveWrapperObject(buildsResponseWithoutEmptyChannel)

	giveMeAName := channelWriteScrapeResultToStorage(storage, scrapeStartTime, buildsWithNoWrappersChannel)

	for _, s := range <-giveMeAName {
		if s == -1 {
			fmt.Print(s)
		}
	}

	// Collect one projectBuildsStrings per project from the channel
	//projectsBuildStrings := make([]string, 0)

	<-buildsResponseAsStringsChannel
	<-buildsResponseWithoutEmptyChannel
	<-buildsWithNoWrappersChannel
	<-giveMeAName

	time.Sleep(10000 * time.Millisecond)
}

func printProgressBar(length int) {

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

func channelRemoveEmptyResults(projectBuildsStrings <-chan string) <-chan string {

	out := make(chan string)

	go func() {
		forwardedBuildStings := 0
		for projectBuildsString := range projectBuildsStrings {
			//fmt.Print("y")
			if projectBuildsString == "{\"count\":0,\"value\":[]}" {
				continue
			}
			out <- projectBuildsString

			forwardedBuildStings++
		}
		close(out)

		//fmt.Sprintf("%d build strings have been forwarded\n", forwardedBuildStings)
		//fmt.Println("")
	}()
	return out

}

func channelRemoveWrapperObject(projectBuildsStrings <-chan string) <-chan string {

	out := make(chan string)

	go func() {
		for projectBuildsString := range projectBuildsStrings {
			projectBuildsString = removeWrapperObject(projectBuildsString)
			out <- projectBuildsString
		}
		close(out)
	}()
	return out

}

func channelWriteScrapeResultToStorage(
	storage Storage,
	scrapeStartTime time.Time,
	projectBuildsStrings <-chan string) <-chan string {

	out := make(chan string)

	go func() {
		fileContent := "["
		for projectBuildsString := range projectBuildsStrings {
			fileContent += projectBuildsString
			out <- projectBuildsString

			//	fmt.Print("█")
		}
		close(out)

		fileContent += "]"

		if len(fileContent) == 2 {
			return
		}

		storage.storeScrapeResult(scrapeStartTime, fileContent)

	}()
	return out

}

type Storage interface {
	storeScrapeResult(timeStamp time.Time, fileContent string)
	getLastScrapeStartTime() time.Time
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
