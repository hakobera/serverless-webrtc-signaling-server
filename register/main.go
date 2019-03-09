package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"

	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type RegisterCommand struct {
	Type     string `json:"type"`
	RoomID   string `json:"room_id"`
	ClientID string `json:"client_id"`
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	now := time.Now().UTC()
	ctx := request.RequestContext
	fmt.Println(ctx.ConnectionID, request.Body)

	cmd := RegisterCommand{}
	err := json.Unmarshal([]byte(request.Body), &cmd)
	if err != nil {
		return common.ErrorResponse(err, 400)
	}

	db := common.DB()
	roomsTable := common.RoomsTable(db)
	var room common.Room

	err = roomsTable.Get("roomId", cmd.RoomID).One(&room)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			room = common.Room{RoomID: cmd.RoomID, Clients: []common.Client{}, Created: now}
		} else {
			return common.ErrorResponse(err, 500)
		}
	}

	result := "accept"
	connectionsTable := common.ConnectionsTable(db)
	conn := common.Connection{ConnectionID: ctx.ConnectionID, RoomID: room.RoomID}

	if len(room.Clients) < 2 {
		client := common.Client{ConnectionID: ctx.ConnectionID, ClientID: cmd.ClientID, Joined: now}
		room.Clients = append(room.Clients, client)
		tx := db.WriteTx()
		tx.Put(roomsTable.Put(room))
		tx.Put(connectionsTable.Put(conn))
		err := tx.Run()
		if err != nil {
			return common.ErrorResponse(err, 500)
		}
	} else {
		result = "reject"
	}

	svc, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	_, err = svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &ctx.ConnectionID,
		Data:         []byte(fmt.Sprintf("{\"type\": \"%s\"}", result)),
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
