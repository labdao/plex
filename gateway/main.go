package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type WorkflowRequest struct {
	Workflow string `json:"workflow"`
}

const IPFSCIDLength = 46 // replace with the correct CID length

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var workflowRequest WorkflowRequest
	err := json.Unmarshal([]byte(request.Body), &workflowRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	if len(workflowRequest.Workflow) != IPFSCIDLength {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid CID length: expected %d, got %d", IPFSCIDLength, len(workflowRequest.Workflow)),
		}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
