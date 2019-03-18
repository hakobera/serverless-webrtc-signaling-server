package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"
)

func ondisconnectHandler(api common.ApiGatewayManagementAPI, db common.DB, connectionID string) error {
	connectionsTable := db.ConnectionsTable()
	roomsTable := db.RoomsTable()

	var conn common.Connection
	var room common.Room
	var err error

	err = connectionsTable.FindOne("connectionId", connectionID, &conn)
	if err != nil {
		return err
	}
	err = connectionsTable.Delete("connectionId", connectionID)
	if err != nil {
		return err
	}

	err = roomsTable.FindOne("roomId", conn.RoomID, &room)
	if err != nil {
		return err
	}

	for _, c := range room.Clients {
		if c.ConnectionID != connectionID {
			// connection might be closed, in that case just ignore error
			api.PostToConnection(c.ConnectionID, `{"type":"close"}`)
		}
	}
	room.Clients = remove(room.Clients, connectionID)
	err = roomsTable.Put(room)
	if err != nil {
		return err
	}

	return nil
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Printf("Disconnected %s\n", ctx.ConnectionID)

	api, err := common.NewApiGatewayManagementApi(ctx.APIID, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	db := common.NewDB(session.New(), aws.NewConfig())
	err = ondisconnectHandler(api, db, ctx.ConnectionID)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func remove(clients []common.Client, connectionID string) []common.Client {
	result := []common.Client{}
	for _, v := range clients {
		if v.ConnectionID != connectionID {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	lambda.Start(handler)
}
