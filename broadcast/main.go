package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"
)

func broadcastHandler(api common.ApiGatewayManagementAPI, db common.DB, connectionID, body string) error {
	connectionsTable := db.ConnectionsTable()
	roomsTable := db.RoomsTable()

	var conn common.Connection
	var room common.Room
	var err error

	err = connectionsTable.FindOne("connectionId", connectionID, &conn)
	if err != nil {
		return err
	}

	err = roomsTable.FindOne("roomId", conn.RoomID, &room)
	if err != nil {
		return err
	}

	var ee error
	for _, c := range room.Clients {
		if c.ConnectionID != connectionID {
			err := api.PostToConnection(c.ConnectionID, body)
			if err != nil {
				ee = err
			}
		}
	}

	if ee != nil {
		return ee
	}

	return nil
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Println(ctx.ConnectionID, ctx.RouteKey, request.Body)

	api, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	db := common.NewDB()
	err = broadcastHandler(api, db, ctx.ConnectionID, request.Body)
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
