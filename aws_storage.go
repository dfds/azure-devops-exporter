package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"sort"
	"strings"
	"time"
)

type awsStorage struct{}

var datalakeS3Bucket = "dfds-datalake"
var folderPrefix = "azure-devops/"

func handleAwsErrors(err error) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			if strings.HasPrefix(awsErr.Message(), "no valid providers in chain.") {
				panic("No AWS credentials found, see following link for more info on how to set AWS credentials:\n https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html")
			}
		} else {
			panicOnError(err)
		}
	}
}

func (awsStorage) getLastScrapeStartTime() time.Time {

	awsSession := getAwsSession()
	svc := s3.New(awsSession)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(datalakeS3Bucket), Prefix: &folderPrefix})

	handleAwsErrors(err)

	var scrapeTimes []string
	for _, file := range resp.Contents {
		fmt.Println("Name:         ", *file.Key)

		scrapeTime := strings.TrimLeft(*file.Key, folderPrefix)
		scrapeTime = strings.TrimLeft(scrapeTime, "azure-devops-builds-")
		scrapeTime = strings.TrimRight(scrapeTime, ".json")

		if scrapeTime == "" {
			continue
		}
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

func getAwsSession() *session.Session {
	var awsConfig *aws.Config
	awsConfig = &aws.Config{
		Region: aws.String("eu-west-1"),
	}

	// The session the S3 Uploader will use
	awsSession := session.Must(session.NewSession(awsConfig))

	return awsSession
}

func (awsStorage) storeScrapeResult(timeStamp time.Time, fileContent string) {

	awsSession := getAwsSession()
	uploader := s3manager.NewUploader(awsSession)

	// Create an uploader with the session and custom options
	// uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
	//     u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
	//     u.Concurrency = 2            // default is 5
	// })

	// Source sample https://paulbradley.org/gos3/

	fileName := "azure-devops-builds-" + timeStamp.Format(time.RFC3339) + ".json"

	buildFileKey := folderPrefix + fileName
	dataToWrite := []byte(fileContent)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(datalakeS3Bucket),
		Key:    aws.String(buildFileKey),
		Body:   bytes.NewReader(dataToWrite),
	})

	handleAwsErrors(err)

	fmt.Printf("file uploaded to, %s\n", result.Location)
}
