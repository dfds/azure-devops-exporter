package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type awsStorage struct{}

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
