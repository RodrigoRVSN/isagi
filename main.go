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
)

type UserRequest struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

func usersHandler(w http.ResponseWriter, request *http.Request) {
	var user UserRequest
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		io.WriteString(w, fmt.Sprintf("Error: %v", err))
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
		io.WriteString(w, fmt.Sprintf("Error loading AWS config: %v", err))
		return
	}

	// awsEndpoint := "http://localhost:4566"
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		// o.BaseEndpoint = aws.String(awsEndpoint)
	})

	bucketName := "rvsnlogs"
	data, _ := json.Marshal(user)
	fileOutput, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(user.ID.String()),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("Error loading AWS config: %v", err))
		return
	}

	fmt.Println(fileOutput.ResultMetadata)
	io.WriteString(w, "Data sent!")
}

func main() {
	http.HandleFunc("/users", usersHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
