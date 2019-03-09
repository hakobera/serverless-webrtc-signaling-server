package common

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func NewApiGatewayManagementApi(domainName, stage string) (*apigatewaymanagementapi.ApiGatewayManagementApi, error) {
	sess, err := session.NewSession(aws.NewConfig().WithEndpoint(fmt.Sprintf("%s/%s", domainName, stage)))
	if err != nil {
		return nil, err
	}

	return apigatewaymanagementapi.New(sess), nil
}

func ErrorResponse(err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	msg := err.Error()
	fmt.Println(statusCode, msg)
	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: statusCode,
	}, err
}
