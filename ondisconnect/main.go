package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"
)

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Printf("Disconnected %s\n", ctx.ConnectionID)

	db := common.DB()
	connectionsTable := common.ConnectionsTable(db)
	roomsTable := common.RoomsTable(db)

	var conn common.Connection
	var room common.Room

	err := connectionsTable.Get("connectionId", ctx.ConnectionID).One(&conn)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}
	err = connectionsTable.Delete("connectionId", ctx.ConnectionID).Run()
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	err = roomsTable.Get("roomId", conn.RoomID).One(&room)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	api, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	for _, c := range room.Clients {
		if c.ConnectionID != ctx.ConnectionID {
			// connection might be closed, in that case just ignore error
			api.PostToConnection(c.ConnectionID, `{"type":"close"}`)
		}
	}
	room.Clients = remove(room.Clients, ctx.ConnectionID)
	err = roomsTable.Put(room).Run()
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
