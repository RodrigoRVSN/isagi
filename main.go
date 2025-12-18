package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const PORT = 8080

var logger *zap.SugaredLogger

type UserRequest struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

func usersHandler(w http.ResponseWriter, request *http.Request) {
	var user UserRequest
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		w.WriteHeader(500)
		logger.Errorf("Error decoding request body %v", err)
		return
	}

	user.ID = uuid.New()
	if user.Timestamp.IsZero() {
		user.Timestamp = time.Now().UTC()
	}

	awsRegion := "us-east-1"
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		w.WriteHeader(500)
		logger.Errorf("Error loading config %v", err)
		return
	}

	// awsEndpoint := "http://localhost:4566"
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		// o.BaseEndpoint = aws.String(awsEndpoint)
	})

	year := user.Timestamp.Year()
	month := user.Timestamp.Month()
	day := user.Timestamp.Day()
	key := fmt.Sprintf("%d/%d/%d/%s", year, month, day, user.ID.String())
	logger.Infof("Saving file %s", key)
	bucketName := "rvsnlogs"
	data, _ := json.Marshal(user)
	fileOutput, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		w.WriteHeader(500)
		logger.Errorf("Error saving file %v", err)
		return
	}

	logger.Infof("File sent! %s", fileOutput.ResultMetadata.Get("key"))
	io.WriteString(w, "Data sent!")
}

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	http.HandleFunc("/users", usersHandler)
	logger = zapLogger.Sugar()
	logger.Infof("Starting server on :%d", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
