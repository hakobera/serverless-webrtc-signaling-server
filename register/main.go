package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/hakobera/serverless-webrtc-signaling-server/common"
)

type RegisterCommand struct {
	Type     string `json:"type"`
	RoomID   string `json:"room_id"`
	ClientID string `json:"client_id"`
}

func registerHandler(api common.ApiGatewayManagementAPI, db common.DB, connectionID, body string) error {
	cmd := RegisterCommand{}
	err := json.Unmarshal([]byte(body), &cmd)
	if err != nil {
		return err
	}

	connectionsTable := db.ConnectionsTable()
	roomsTable := db.RoomsTable()
	now := time.Now().UTC()

	var room common.Room
	err = roomsTable.FindOne("roomId", cmd.RoomID, &room)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			room = common.Room{RoomID: cmd.RoomID, Clients: []common.Client{}, Created: now}
		} else {
			return err
		}
	}

	conn := common.Connection{ConnectionID: connectionID, RoomID: room.RoomID}
	result := "accept"

	if len(room.Clients) < 2 {
		client := common.Client{ConnectionID: connectionID, ClientID: cmd.ClientID, Joined: now}
		room.Clients = append(room.Clients, client)
		err := db.TxPut(
			&common.TableItem{Table: roomsTable, Item: room},
			&common.TableItem{Table: connectionsTable, Item: conn},
		)
		if err != nil {
			return err
		}
	} else {
		result = "reject"
	}

	return api.PostToConnection(connectionID, fmt.Sprintf("{\"type\": \"%s\"}", result))
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := request.RequestContext
	fmt.Println(ctx.ConnectionID, request.Body)

	api, err := common.NewApiGatewayManagementApi(ctx.DomainName, ctx.Stage)
	if err != nil {
		return common.ErrorResponse(err, 500)
	}

	db := common.NewDB(session.New(), aws.NewConfig())
	err = registerHandler(api, db, ctx.ConnectionID, request.Body)
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
