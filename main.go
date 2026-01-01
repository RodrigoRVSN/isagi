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
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
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
	request.Header.Set("Content-Type", "application/json")
	awsRegion := "us-east-1"
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		w.WriteHeader(500)
		logger.Errorf("Error loading config %v", err)
		return
	}

	if request.Method == "POST" {
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
		bucketName := "rvsn-logs"
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

		// TODO: save parquet in AWS (other bucket)
		// parquetErr := parquet.WriteFile("file.parquet", []UserRequest{user})
		// if parquetErr != nil {
		// 	w.WriteHeader(500)
		// 	logger.Errorf("Error saving file %v", err)
		// 	return
		// }

		logger.Infof("File sent! %s", fileOutput.ResultMetadata.Get("key"))
		io.WriteString(w, "Data sent!")
	}

	if request.Method == "GET" {
		awsClient := athena.NewFromConfig(awsCfg)
		output, err := awsClient.StartQueryExecution(
			context.TODO(),
			&athena.StartQueryExecutionInput{
				WorkGroup:   aws.String("workgroup_name"),
				QueryString: aws.String("SELECT * FROM table_name"),
				QueryExecutionContext: &types.QueryExecutionContext{
					Database: aws.String("logs"),
				},
			},
		)

		if err != nil {
			w.WriteHeader(500)
			logger.Errorf("Error loading config %v", err)
			return
		}

		var status types.QueryExecutionState
		for status != types.QueryExecutionStateSucceeded {
			execution, err := awsClient.GetQueryExecution(context.TODO(), &athena.GetQueryExecutionInput{QueryExecutionId: output.QueryExecutionId})
			if err != nil {
				w.WriteHeader(500)
				logger.Errorf("Error getting query results %v", err)
				return
			}
			status = execution.QueryExecution.Status.State
		}

		var users []UserRequest
		results, err := awsClient.GetQueryResults(context.TODO(), &athena.GetQueryResultsInput{QueryExecutionId: output.QueryExecutionId})

		if err != nil {
			w.WriteHeader(500)
			logger.Errorf("Error loading config %v", err)
			return
		}

		for _, row := range results.ResultSet.Rows {
			userId, _ := uuid.Parse(*row.Data[0].VarCharValue)
			users = append(users, UserRequest{ID: userId})
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	http.HandleFunc("/users", usersHandler)
	logger = zapLogger.Sugar()
	logger.Infof("Starting server on :%d", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
