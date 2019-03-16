package common

import (
	"fmt"

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

func NewApiGatewayManagementApi(domainName, stage string) (ApiGatewayManagementAPI, error) {
	sess, err := session.NewSession(aws.NewConfig().WithEndpoint(fmt.Sprintf("%s/%s", domainName, stage)))
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
