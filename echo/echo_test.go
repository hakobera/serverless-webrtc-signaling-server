package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"

	mock "github.com/hakobera/serverless-webrtc-signaling-server/mock_common"
)

func TestSample1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connectionID := "Conn1"
	body := "message"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	api.EXPECT().PostToConnection(connectionID, body).Return(nil)

	request := events.APIGatewayWebsocketProxyRequest{
		Body: body,
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			ConnectionID: connectionID,
		},
	}
	echoHandler(request, api)
}
