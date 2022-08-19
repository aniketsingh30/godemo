package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

const ID = "id"
const RESULT = "result"
const STATUS_CODE = "statusCode"
const SUCCESS_MSG = "successMessage"
const FAILLED_MSG = "failedMessage"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetFinalResponse(id uuid.UUID, result interface{}, status int, succeeMessage string, err error) (*events.APIGatewayProxyResponse, error) {
	finalResponse := make(map[string]interface{})
	finalResponse[ID] = id
	finalResponse[RESULT] = result
	finalResponse[STATUS_CODE] = status
	if succeeMessage != "" {
		finalResponse[SUCCESS_MSG] = succeeMessage
	}
	if err != nil {
		finalResponse[FAILLED_MSG] = ErrorBody{aws.String(err.Error())}
	}
	return apiResponse(status, finalResponse), nil
}

func GetErrorResponse(id uuid.UUID, result interface{}, status int, succeeMessage string, err error) (*events.APIGatewayProxyResponse, error) {
	finalResponse := make(map[string]interface{})

	finalResponse[ID] = id
	finalResponse[RESULT] = result
	finalResponse[STATUS_CODE] = status
	if succeeMessage != "" {
		finalResponse[SUCCESS_MSG] = succeeMessage
	}
	if err != nil {
		finalResponse[FAILLED_MSG] = ErrorBody{aws.String(err.Error())}
	}
	return apiResponse(status, finalResponse), err
}

func apiResponse(status int, body interface{}) *events.APIGatewayProxyResponse {
	headerMap := map[string]string{"Content-Type": "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "OPTIONS,POST,PATCH,GET,DELETE,PUT"}

	resp := events.APIGatewayProxyResponse{Headers: headerMap}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp
}
