package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/guregu/dynamo"

	"github.com/hakobera/serverless-webrtc-signaling-server/common"
	mock "github.com/hakobera/serverless-webrtc-signaling-server/mock_common"
)

func TestRegisterHandlerAccept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	roomID := "Room1"
	body := fmt.Sprintf(`{"type": "register", "room_id": "%s", "client_id": "client1"}`, roomID)
	acceptMessage := `{"type": "accept"}`
	now := time.Now().UTC()

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, acceptMessage).Return(nil)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)
	db.EXPECT().TxPut(
		common.TableItem{
			Table: roomsTable,
			Item: common.Room{
				RoomID: roomID,
				Clients: []common.Client{
					common.Client{
						ConnectionID: myConnID,
						ClientID:     "client1",
						Joined:       now,
					},
				},
				Created: now,
			},
		},
		common.TableItem{
			Table: connectionsTable,
			Item:  common.Connection{ConnectionID: myConnID, RoomID: roomID},
		},
	).Return(nil)

	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			return dynamo.ErrNotFound
		})
	registerHandler(api, db, myConnID, body, now)
}

func TestRegisterHandlerReject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	otherConnID1 := "Other1"
	otherConnID2 := "Other2"
	roomID := "Room1"
	body := fmt.Sprintf(`{"type": "register", "room_id": "%s", "client_id": "client1"}`, roomID)
	acceptMessage := `{"type": "reject"}`

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, acceptMessage).Return(nil)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)
	db.EXPECT().TxPut(gomock.Any()).Return(nil).Times(0)

	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			room := out.(*common.Room)
			room.RoomID = roomID
			room.Clients = []common.Client{
				common.Client{
					ConnectionID: otherConnID1,
				},
				common.Client{
					ConnectionID: otherConnID2,
				},
			}
			return nil
		})

	registerHandler(api, db, myConnID, body, time.Now().UTC())
}

func TestRegisterHandlerRequestBodyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	roomID := "Room1"
	body := "invalid body"
	acceptMessage := `{"type": "reject"}`

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, acceptMessage).Return(nil).Times(0)

	db.EXPECT().ConnectionsTable().Return(connectionsTable).Times(0)
	db.EXPECT().RoomsTable().Return(roomsTable).Times(0)
	db.EXPECT().TxPut(gomock.Any()).Return(nil).Times(0)

	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).Return(nil).Times(0)

	err := registerHandler(api, db, myConnID, body, time.Now().UTC())
	if err.Error() != "invalid character 'i' looking for beginning of value" {
		t.Errorf(err.Error())
	}
}
