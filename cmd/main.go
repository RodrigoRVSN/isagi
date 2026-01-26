package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rodrigorvsn/isagi/cmd/create_handler"
	"github.com/rodrigorvsn/isagi/cmd/get_handler"
	"github.com/rodrigorvsn/isagi/cmd/lib/logger"
)

const PORT = 8080

func usersHandler(w http.ResponseWriter, request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	awsRegion := "us-east-1"
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Errorf("Error loading config %v", err)
		return
	}

	switch request.Method {
	case http.MethodPost:
		create_handler.Handle(request, w, awsCfg)
	case http.MethodGet:
		get_handler.Handle(w, awsCfg)
	}
}

func main() {
	logger := logger.Init()
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	logger.Infof("Starting server on :%d", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
