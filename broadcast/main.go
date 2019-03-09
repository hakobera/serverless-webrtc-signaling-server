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
	fmt.Println(ctx.ConnectionID, ctx.RouteKey, request.Body)

	db := common.DB()
	connectionsTable := common.ConnectionsTable(db)
	roomsTable := common.RoomsTable(db)

	var conn common.Connection
	var room common.Room

	err := connectionsTable.Get("connectionId", ctx.ConnectionID).One(&conn)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	err = roomsTable.Get("roomId", conn.RoomID).One(&room)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	svc, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	//TODO: improve error handling
	var ee error
	for _, c := range room.Clients {
		if c.ConnectionID != ctx.ConnectionID {
			_, err = svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &c.ConnectionID,
				Data:         []byte(request.Body),
			})

			if err != nil {
				ee = err
			}
		}
	}

	if ee != nil {
		return common.ErrorResponse(ee, 500)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Data sent.",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
