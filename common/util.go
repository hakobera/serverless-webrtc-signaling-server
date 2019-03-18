package common

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ApiGatewayManagementAPIImpl struct {
	client *apigatewaymanagementapi.ApiGatewayManagementApi
}

func (a *ApiGatewayManagementAPIImpl) PostToConnection(connectionID, body string) error {
	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionID,
		Data:         []byte(body),
	}
	_, err := a.client.PostToConnection(input)
	return err
}

func NewApiGatewayManagementApi(APIID, stage string) (ApiGatewayManagementAPI, error) {
	region := os.Getenv("AWS_REGION")
	endpoint := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s", APIID, region, stage)
	sess, err := session.NewSession(aws.NewConfig().WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	api := ApiGatewayManagementAPIImpl{
		client: apigatewaymanagementapi.New(sess),
	}
	return &api, nil
}

func ErrorResponse(err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	msg := err.Error()
	fmt.Println(statusCode, msg)
	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: statusCode,
	}, err
}
