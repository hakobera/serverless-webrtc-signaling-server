package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/hakobera/serverless-webrtc-signaling-server/common"
	mock "github.com/hakobera/serverless-webrtc-signaling-server/mock_common"
)

func TestBroadcastHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	otherConnID := "Other1"
	body := "message"
	roomID := "Room1"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, body).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, body).Return(nil).Times(1)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			conn := out.(*common.Connection)
			conn.ConnectionID = myConnID
			conn.RoomID = roomID
			return nil
		})
	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			room := out.(*common.Room)
			room.RoomID = roomID
			room.Clients = []common.Client{
				common.Client{
					ConnectionID: myConnID,
				},
				common.Client{
					ConnectionID: otherConnID,
				},
			}
			return nil
		})

	broadcastHandler(api, db, myConnID, body)
}

func TestBroadcastHandlerConnectionNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	body := "message"
	otherConnID := "Other1"
	expectedError := "connection not found"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, body).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, body).Return(nil).Times(0)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			return fmt.Errorf(expectedError)
		})

	err := broadcastHandler(api, db, myConnID, body)
	if err.Error() != expectedError {
		t.Errorf("expect %s, but actual %s", expectedError, err.Error())
	}
}

func TestBroadcastHandlerRoomNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	myConnID := "Conn1"
	otherConnID := "Other1"
	body := "message"
	roomID := "Room1"
	expectedError := "room not found"

	api := mock.NewMockApiGatewayManagementAPI(ctrl)
	db := mock.NewMockDB(ctrl)
	connectionsTable := mock.NewMockTable(ctrl)
	roomsTable := mock.NewMockTable(ctrl)

	api.EXPECT().PostToConnection(myConnID, body).Return(nil).Times(0)
	api.EXPECT().PostToConnection(otherConnID, body).Return(nil).Times(0)

	db.EXPECT().ConnectionsTable().Return(connectionsTable)
	db.EXPECT().RoomsTable().Return(roomsTable)

	connectionsTable.EXPECT().FindOne("connectionId", myConnID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			conn := out.(*common.Connection)
			conn.ConnectionID = myConnID
			conn.RoomID = roomID
			return nil
		})
	roomsTable.EXPECT().FindOne("roomId", roomID, gomock.Any()).DoAndReturn(
		func(column string, key, out interface{}) error {
			return fmt.Errorf(expectedError)
		})

	err := broadcastHandler(api, db, myConnID, body)
	if err.Error() != expectedError {
		t.Errorf("expect %s, but actual %s", expectedError, err.Error())
	}
}
