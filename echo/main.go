package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"

	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Println(ctx.ConnectionID, request.Body)

	svc, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	_, err = svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &ctx.ConnectionID,
		Data:         []byte(request.Body),
	})

	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Data sent.",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
