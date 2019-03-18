package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"
)

func echoHandler(request events.APIGatewayWebsocketProxyRequest, api common.ApiGatewayManagementAPI) error {
	ctx := request.RequestContext
	return api.PostToConnection(ctx.ConnectionID, request.Body)
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Printf("[echo] connectionID=%s, body=%s\n", ctx.ConnectionID, request.Body)

	api, err := common.NewApiGatewayManagementApi(ctx.APIID, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	err = echoHandler(request, api)
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
