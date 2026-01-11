package get_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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

	results, err := awsClient.GetQueryResults(context.TODO(), &athena.GetQueryResultsInput{QueryExecutionId: output.QueryExecutionId})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Errorf("Error loading config %v", err)
		return
	}

	var columnIndex map[string]int
	columnIndex = make(map[string]int)
	for i, col := range results.ResultSet.ResultSetMetadata.ColumnInfo {
		columnIndex[*col.Name] = i
	}

	var users []models.User
	for indexRow, row := range results.ResultSet.Rows {
		if indexRow == 0 { // ignore headers
			continue
		}
		var user models.User

		if idx, ok := columnIndex["id"]; ok {
			if row.Data[idx].VarCharValue != nil {
				if id, err := uuid.Parse(*row.Data[idx].VarCharValue); err == nil {
					user.ID = id
				}
			}
		}
		if idx, ok := columnIndex["username"]; ok {
			if row.Data[idx].VarCharValue != nil {
				user.Username = *row.Data[idx].VarCharValue
			}
		}
		if idx, ok := columnIndex["timestamp"]; ok {
			if row.Data[idx].VarCharValue != nil {
				if time, err := time.Parse(time.RFC3339, *row.Data[idx].VarCharValue); err == nil {
					user.Timestamp = time
				}
			}
		}

		users = append(users, user)

	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
