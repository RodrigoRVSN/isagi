package create_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rodrigorvsn/isagi/cmd/lib/logger"
	"github.com/rodrigorvsn/isagi/cmd/models"
)

func Handle(request *http.Request, writer http.ResponseWriter, awsCfg aws.Config) {
	var user models.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Logger().Errorf("Error decoding request body %v", err)
		return
	}

	user.ID = uuid.New()
	if user.Timestamp.IsZero() {
		user.Timestamp = time.Now().UTC()
	}

	// awsEndpoint := "http://localhost:4566"
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		// o.BaseEndpoint = aws.String(awsEndpoint)
	})

	year := user.Timestamp.Year()
	month := user.Timestamp.Month()
	day := user.Timestamp.Day()
	key := fmt.Sprintf("%d/%d/%d/%s.json", year, month, day, user.ID.String())
	logger.Logger().Infof("Saving file %s", key)
	bucketName := "rvsn-logs"
	data, _ := json.Marshal(user)
	fileOutput, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		writer.WriteHeader(500)
		logger.Logger().Errorf("Error saving file %v", err)
		return
	}

	logger.Logger().Infof("File sent! %s", fileOutput.ResultMetadata.Get("key"))
	writer.WriteHeader(http.StatusCreated)
	io.WriteString(writer, "Data sent!")
}
