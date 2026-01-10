package get_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/google/uuid"
	"github.com/rodrigorvsn/isagi/cmd/lib/logger"
	"github.com/rodrigorvsn/isagi/cmd/models"
)

func Handle(w http.ResponseWriter, awsCfg aws.Config) {
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
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Errorf("Error loading config %v", err)
		return
	}

	var status types.QueryExecutionState
	for status != types.QueryExecutionStateSucceeded {
		execution, err := awsClient.GetQueryExecution(context.TODO(), &athena.GetQueryExecutionInput{QueryExecutionId: output.QueryExecutionId})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Logger().Errorf("Error getting query results %v", err)
			return
		}
		status = execution.QueryExecution.Status.State
	}

	var users []models.User
	results, err := awsClient.GetQueryResults(context.TODO(), &athena.GetQueryResultsInput{QueryExecutionId: output.QueryExecutionId})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Errorf("Error loading config %v", err)
		return
	}

	for _, row := range results.ResultSet.Rows {
		userId, _ := uuid.Parse(*row.Data[0].VarCharValue)
		users = append(users, models.User{ID: userId})
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
